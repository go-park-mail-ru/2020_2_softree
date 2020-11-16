package config

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

// GlobalConfig instance of application configuration
var GlobalConfig = viper.New()

// ParseConfig push application settings in global store
func ParseConfig(filename string, defaults map[string]interface{}) error {
	for key, value := range defaults {
		GlobalConfig.SetDefault(key, value)
	}

	fullpath, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	GlobalConfig.SetConfigType("yaml")
	GlobalConfig.SetConfigFile(fullpath)
	GlobalConfig.AutomaticEnv()
	if err := GlobalConfig.ReadInConfig(); err != nil {
		return err
	}
	createURLS()
	return nil
}

func createURLS() {
	GlobalConfig.Set("redis.sessionURL", fmt.Sprintf("redis://%s:%s:%d%s",
		GlobalConfig.GetString("redis.user"),
		GlobalConfig.GetString("redis.host"),
		GlobalConfig.GetInt("redis.port"),
		GlobalConfig.GetString("redis.sessionPath")))

	GlobalConfig.Set("redis.currencyURL", fmt.Sprintf("redis://%s:%s:%d%s",
		GlobalConfig.GetString("redis.user"),
		GlobalConfig.GetString("redis.host"),
		GlobalConfig.GetInt("redis.port"),
		GlobalConfig.GetString("redis.currencyPath")))

	GlobalConfig.Set("postgres.URL", fmt.Sprintf(
		"host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		GlobalConfig.GetString("postgres.host"),
		GlobalConfig.GetInt("postgres.port"),
		GlobalConfig.GetString("postgres.user"),
		GlobalConfig.GetString("postgres.password"),
		GlobalConfig.GetString("postgres.db")))
}
