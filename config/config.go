// Package config holds the configuration structs for your application or service
// You would likely define either a single file per config struct (optimal) or place
// all config structs in this file (suboptimal but not incorrect)
package config

import "time"

// Config is a collection of configuration variables and pointers to dependencies.
type Config struct {
	HTTPPort int `env:"HTTP_PORT" default:"1234"`
	// TODO this needs to pivot to be multiple versions
	HTTPApiVersion int `env:"HTTP_API_VERSION" default:"1"`

	PBPort int `env:"PB_PORT" default:"1234"`
	// TODO this needs to pivot to be multiple versions
	PBApiVersion int    `env:"PB_API_VERSION" default:"1"`
	DBConnString string `env:"DB_CONN_STRING"`

	StatsdPrefix   string        `env:"STATSD_PREFIX" default:"xxx"`
	StatsdAddress  string        `env:"STATD_ADDRESS" default:"127.0.0.1"`
	StatsdInterval time.Duration `env:"STATSD_INTERVAL" default:"10"`
}
