package config

import (
	"github.com/roguexray007/loan-app/pkg/db"
	"github.com/roguexray007/loan-app/pkg/mutex"
)

// Config holds all the config required for the application
type Config struct {
	App   App
	Db    db.ConnectionsConfig `mapstructure:"db"`
	Redis Redis
	Mutex mutex.Config
}

// App contains application-specific config values
type App struct {
	Env             string
	ServiceName     string
	Hostname        string
	Port            string
	MetricPort      string
	ShutdownTimeout int
	ShutdownDelay   int
	GitCommitHash   string
	Debug           bool
}

type Redis struct {
	Host     string
	Port     int32
	Database int32
	Password string
	Mode     string
	Dialect  string
}
