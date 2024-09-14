package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	HTTP struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}
	Postgres struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DB       string `yaml:"db"`
	}
}

// NewConfig - создание новой конфигурации
func NewConfig(path string) (*Config, error) {
	cfg := &Config{}
	err := LoadConfig(path, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// LoadConfig - загрузка конфигурации
func LoadConfig(path string, cfg *Config) error {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return err
	}

	return nil
}
