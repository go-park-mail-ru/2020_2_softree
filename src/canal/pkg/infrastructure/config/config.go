package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"path/filepath"

	"github.com/spf13/viper"
)

// ParseConfig set defaults and read from config file
func ParseConfig(filename string, defaults map[string]interface{}) error {
	for key, value := range defaults {
		viper.SetDefault(key, value)
	}

	fullpath, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	viper.SetConfigType("yaml")
	viper.SetConfigFile(fullpath)
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	viper.WatchConfig()
	viper.OnConfigChange(
		func(event fsnotify.Event) {
			fmt.Println("changed with event: ", event.Name)
		},
	)

	createURLS()
	return nil
}

func createURLS() {
	viper.Set("redis.sessionURL", fmt.Sprintf("redis://%s:%s:%d%s",
		viper.GetString("redis.user"),
		viper.GetString("redis.host"),
		viper.GetInt("redis.port"),
		viper.GetString("redis.sessionPath")))

	viper.Set("redis.currencyURL", fmt.Sprintf("redis://%s:%s:%d%s",
		viper.GetString("redis.user"),
		viper.GetString("redis.host"),
		viper.GetInt("redis.port"),
		viper.GetString("redis.currencyPath")))

	viper.Set("postgres.URL", fmt.Sprintf(
		"host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		viper.GetString("postgres.host"),
		viper.GetInt("postgres.port"),
		viper.GetString("postgres.user"),
		viper.GetString("postgres.password"),
		viper.GetString("postgres.db")))
}
