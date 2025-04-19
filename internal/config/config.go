package config

import (
	"fmt"
)

type CliConfig struct {
	Host string
	Port string
}

func (cfg CliConfig) Url() string {
	return fmt.Sprintf("http://%s:%s/", cfg.Host, cfg.Port)
}
