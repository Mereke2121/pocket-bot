package main

import (
	"log"

	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatal("initialize config")
	}
}

func initConfig() error {
	viper.SetConfigName("conf")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	return viper.ReadInConfig()
}
