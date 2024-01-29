package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ServerAddress int `yaml:"server_address"`
}

func LoadConfig(path string, config *Config) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return err
	}

	return nil
}
