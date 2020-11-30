package config

import (
	"path/filepath"

	"github.com/spf13/viper"
)

// ParseConfig set defaults and read from config file
func ParseConfig(filename string, defaults map[string]interface{}) error {
	for key, value := range defaults {
		viper.SetDefault(key, value)
	}

	fullPath, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	viper.SetConfigType("yaml")
	viper.SetConfigFile(fullPath)
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
