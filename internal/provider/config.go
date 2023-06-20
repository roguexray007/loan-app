package provider

import (
	"fmt"
	"os"
	"sync"

	"github.com/roguexray007/loan-app/pkg/container"

	"github.com/roguexray007/loan-app/internal/config"
	"github.com/roguexray007/loan-app/internal/trace"

	config_reader "github.com/roguexray007/loan-app/pkg/config"
)

type appConfig struct {
	sync.Mutex
}

// Config : key which holds the config in container
const Config = "config"

func init() {
	dep.Register(Config, &appConfig{})
}

// Build will build the new instance of config
// by loading it from the config file with respect to the application env
func (cc *appConfig) Build(c container.IContainer) (container.IDependency, error) {
	cc.Lock()
	defer cc.Unlock()

	var cnf *config.Config

	err := config_reader.NewDefaultConfig().Load(c.GetEnv(), &cnf)
	if err != nil {
		fmt.Printf(trace.ConfigLoadingError+" %v", err)
		return nil, err
	}

	// set static dynamic config values
	cnf.App.Env = c.GetEnv()
	cnf.App.GitCommitHash = os.Getenv("GIT_COMMIT_HASH")

	// register the data in Container
	c.Put(Config, cnf)

	return cnf, nil
}
