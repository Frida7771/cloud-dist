package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Load reads YAML/TOML/JSON config into Config struct.
func Load(path string) (Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("unmarshal config: %w", err)
	}

	return cfg, nil
}
