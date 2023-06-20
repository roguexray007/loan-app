package provider

import (
	"context"
	"os"

	"github.com/fatih/structs"

	"github.com/roguexray007/loan-app/internal/config"
	"github.com/roguexray007/loan-app/internal/constants"
	"github.com/roguexray007/loan-app/pkg/container"
	"github.com/roguexray007/loan-app/pkg/db"
)

// Fields type, used to pass to `WithFields`.
type Fields map[string]interface{}

func init() {
	container.
		Init(GetContext(context.Background()), getEnv(), GetManager())
}

// MustResolve will resolve the dependency of given list
// failed to resolve then throws panic
func MustResolve(list []string) {
	container.GetContainer().MustResolve(list)
}

// GetContext will give the context of the application
func GetContext(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.WithValue(context.Background(), constants.ContextKeyApp, structs.Map(GetConfig(ctx).App))
	}

	return ctx
}

// GetConfig will resolve the load the config if not already loaded and provide the config
func GetConfig(_ context.Context) *config.Config {
	if value, ok := container.Get(Config); ok {
		return value.(*config.Config)
	}
	return nil
}

// GetDatabase will resolve the database provider and provide the db instance
func GetDatabase(_ context.Context) *db.Connections {
	if value, ok := container.Get(Database); ok {
		return value.(*db.Connections)
	}

	return nil
}

func getEnv() string {
	// Fetch env for bootstrapping
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	return env
}
