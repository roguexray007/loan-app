package db

import (
	"context"
	"database/sql"
	"reflect"

	"github.com/jinzhu/gorm"
)

const (
	Master  = "master"
	Replica = "replica"
)

type ConnectionsConfig struct {
	Master  Config `mapstructure:"master"`
	Replica Config `mapstructure:"replica"`
}

type Connections struct {
	master  *DB
	replica *DB
}

func NewConnections(c *ConnectionsConfig, args ...interface{}) (*Connections, error) {

	dbConnections := &Connections{}

	var err error

	if len(args) > 1 {
		return nil, ErrorInvalidArguments
	}

	if len(args) == 1 {
		driverMap, ok := args[0].(map[string]*sql.DB)

		if !ok {
			return nil, ErrorInvalidArguments
		}

		dbConnections.master, err = NewDb(&c.Master, driverMap[Master])

		if err != nil {
			return nil, err
		}

		dbConnections.replica, err = NewDb(&c.Replica, driverMap[Replica])

		if err != nil {
			return nil, err
		}

		return dbConnections, nil
	}

	dbConnections.master, err = NewDb(&c.Master)

	if err != nil {
		return nil, err
	}

	dbConnections.replica, err = NewDb(&c.Replica)

	if err != nil {
		return nil, err
	}

	return dbConnections, nil
}

func (c Connections) GetDbByConnectionName(name string) *DB {
	switch name {
	case Master:
		return c.master

	case Replica:
		return c.replica

	default:
		return nil
	}
}

func (c Connections) GetConnection(ctx context.Context, args ...interface{}) *gorm.DB {
	var database *DB
	var connection string
	forceInstance := true

	if len(args) == 0 {
		connection = GetDefaultConnectionFromCtx(ctx)

		forceInstance = false
	} else {
		connection, _ = args[0].(string)

		if len(args) > 1 {
			forceInstance, _ = args[1].(bool)
		}
	}

	database = c.GetDbByConnectionName(connection)

	if database == nil {
		database = c.master
	}

	if forceInstance {
		return database.ForceInstance(ctx)
	}

	return database.Instance(ctx)
}

func (c Connections) Destroy() error {
	var err error

	v := reflect.ValueOf(c)

	for i := 0; i < v.NumField(); i++ {
		db := c.GetDbByConnectionName(v.Type().Field(i).Name)

		if db == nil {
			continue
		}

		err1 := db.instance.Close()

		if err1 != nil {
			err = err1
		}
	}

	return err
}

// Ping checks for master db by default. Ping for other DBs can be checked by passing the connection name in args.
func (c Connections) Ping(ctx context.Context, args ...interface{}) error {
	instance := c.GetConnection(ctx, args...)

	return instance.DB().Ping()
}

func GetDefaultConnectionFromCtx(ctx context.Context) string {
	if defaultConnection, ok := ctx.Value(ContextKeyDatabaseConnection).(string); ok {
		return defaultConnection
	}

	return Master
}
