package ai

import (
	"github.com/kelseyhightower/envconfig"
)

func CreateConfig() (*Aiconfig, error) {
	var cfg Aiconfig
	err := envconfig.Process("AI", &cfg)
	if err != nil {
		return nil, err
	}
	if cfg.Provider != "google" {
		return nil, envconfig.ErrInvalidSpecification
	}

	return &cfg, nil
}

type Aiconfig struct {
	Provider    string  `default:"google"`
	ApiKey      string  `default:"app"`
	Model       string  `default:"gemini-2.0-flash"`
	Temperature float64 `default:"0.2"`
}
