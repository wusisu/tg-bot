package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetDefault("DbPath", "./db.sqlite")
	viper.SetDefault("DownloadDir", "/data/")
	viper.SetDefault("BotToken", "{BotID}:{BotKey}")
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s ", err))
	}
}
