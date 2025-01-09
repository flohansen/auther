package application

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

func (c PostgresConfig) Dsn() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.Username, c.Password, c.Host, c.Port, c.Database)
}

type Config struct {
	Postgres PostgresConfig `yaml:"postgres"`
}

func LoadConfig(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("could not open file: %w", err)
	}

	var config Config
	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		return Config{}, fmt.Errorf("could not decode config: %w", err)
	}

	return config, nil
}
