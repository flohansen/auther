package application

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type SSLMode int

const (
	SSLModeRequire SSLMode = iota
	SSLModeDisable
	SSLModeVerifyCA
	SSLModeVerifyFull
)

func (p SSLMode) String() string {
	switch p {
	case SSLModeDisable:
		return "disable"
	case SSLModeRequire:
		return "require"
	case SSLModeVerifyCA:
		return "verify-ca"
	case SSLModeVerifyFull:
		return "verify-full"
	default:
		return "unknown"
	}
}

func (p *SSLMode) UnmarshalYAML(value *yaml.Node) error {
	var val string
	if err := value.Decode(&val); err != nil {
		return err
	}

	switch val {
	case "require":
		*p = SSLModeRequire
	case "disable":
		*p = SSLModeDisable
	case "verify-ca":
		*p = SSLModeVerifyCA
	case "verify-full":
		*p = SSLModeVerifyFull
	default:
		return fmt.Errorf("unknown value for sslmode: %s, allowed values: ['require', 'disable', 'verify-ca', 'verify-full']", val)
	}

	return nil
}

type PostgresConfig struct {
	Host     string  `yaml:"host"`
	Port     int     `yaml:"port"`
	Username string  `yaml:"username"`
	Password string  `yaml:"password"`
	Database string  `yaml:"database"`
	SSLMode  SSLMode `yaml:"sslmode"`
}

func (c PostgresConfig) Dsn() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Username, c.Password, c.Host, c.Port, c.Database, c.SSLMode)
}

type OIDCConfig struct {
	PrivateKey string `yaml:"privateKey"`
	PublicKey  string `yaml:"publicKey"`
}

type Config struct {
	Postgres PostgresConfig `yaml:"postgres"`
	OIDC     OIDCConfig     `yaml:"oidc"`
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
