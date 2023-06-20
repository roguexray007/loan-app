package config

import "github.com/roguexray007/loan-app/pkg/db"

// Config holds all the config required for the application
type Config struct {
	App App
	Db  db.ConnectionsConfig `mapstructure:"db"`
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
