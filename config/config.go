// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period time.Duration `config:"period"`
	InsecureSkipVerify bool `config:"InsecureSkipVerify"`
	Host    string        `config:"host"`
	Include []string      `config:"include"`
}

var DefaultConfig = Config{
	Period: 1 * time.Second,
	InsecureSkipVerify: true,
	Include: []string{"system.load.average.1m","tomcat.global.error"},
}
