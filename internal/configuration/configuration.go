// Package configuration provides the Viper-based configuration infrastructure for prst.
package configuration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/spf13/viper"
)

// NewViper creates and configures a Viper instance for prst.
func NewViper() (*viper.Viper, error) {
	v := viper.New()

	v.SetEnvPrefix("PRST")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	configPath, err := xdg.ConfigFile("prst/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to resolve config path: %w", err)
	}

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(filepath.Dir(configPath))

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	return v, nil
}

// EnsureConfigDir creates the prst configuration directory if it does not exist.
func EnsureConfigDir() (string, error) {
	configPath, err := xdg.ConfigFile("prst/config.yaml")
	if err != nil {
		return "", fmt.Errorf("failed to resolve config path: %w", err)
	}

	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return configDir, nil
}
