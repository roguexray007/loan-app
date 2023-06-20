package provider

import (
	"sync"

	"github.com/roguexray007/loan-app/internal/config"
	"github.com/roguexray007/loan-app/pkg/container"
	"github.com/roguexray007/loan-app/pkg/db"
)

type database struct {
	sync.Mutex
}

// Database: key which holds the database instance in container
const Database = "database"

func init() {
	dep.Register(Database, &database{})
}

// Build will build the new instance of database
// it'll use the application config to initialize the database
func (dbc *database) Build(c container.IContainer) (container.IDependency, error) {
	dbc.Lock()
	defer dbc.Unlock()

	cfg, _ := c.Get(Config)
	dbCfg := cfg.(*config.Config)

	DB, err := db.NewConnections(&dbCfg.Db)

	if err != nil {
		return nil, err
	}

	c.Put(Database, DB)

	return DB, nil
}

// Destroy will close the connection for database
func (dbc *database) Destroy(c container.IContainer) {
	val, ok := c.Get(Database)
	if !ok {
		return
	}

	dbConnections := val.(*db.Connections)

	_ = dbConnections.Destroy
}
