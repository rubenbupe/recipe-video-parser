package storage

import (
	"github.com/kelseyhightower/envconfig"
)

func CreateConfig() (*Dbconfig, error) {
	var cfg Dbconfig
	err := envconfig.Process("DB", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
