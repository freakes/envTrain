package config

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   Server   `mapstructure:"server"`
	Log      Log      `mapstructure:"log"`
	Database Database `mapstructure:"database"`
}
type Server struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

func (s Server) Address() string { return net.JoinHostPort(s.Host, strconv.Itoa(s.Port)) }

type Log struct {
	Level string `mapstructure:"level"`
	JSON  bool   `mapstructure:"json"`
}
type Database struct {
	Driver string `mapstructure:"driver"`
	DSN    string `mapstructure:"dsn"`
}

// Load uses defaults, an optional explicitly selected JSON file, then APP_* overrides.
func Load(path string) (Config, error) {
	v := viper.New()

	defaults := map[string]any{
		"server.host": "0.0.0.0", "server.port": 8080,
		"server.read_timeout": "10s", "server.write_timeout": "15s", "server.shutdown_timeout": "10s",
		"log.level": "info", "log.json": true, "database.driver": "memory", "database.dsn": "",
	}
	for key, value := range defaults {
		v.SetDefault(key, value)
	}

	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if path != "" {
		v.SetConfigFile(path)

		if err := v.ReadInConfig(); err != nil {
			return Config{}, fmt.Errorf("read config: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return cfg, fmt.Errorf("decode config: %w", err)
	}

	if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
		return cfg, fmt.Errorf("server.port must be between 1 and 65535")
	}

	if cfg.Server.ReadTimeout <= 0 || cfg.Server.WriteTimeout <= 0 || cfg.Server.ShutdownTimeout <= 0 {
		return cfg, fmt.Errorf("server timeouts must be positive")
	}

	if cfg.Database.Driver != "memory" && cfg.Database.Driver != "postgres" {
		return cfg, fmt.Errorf("database.driver must be memory or postgres")
	}

	if cfg.Database.Driver == "postgres" && cfg.Database.DSN == "" {
		return cfg, fmt.Errorf("database.dsn is required for postgres")
	}

	return cfg, nil
}
