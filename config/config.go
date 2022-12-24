package config

import "github.com/spf13/viper"

func InitConfig() error {
	viper.SetConfigName("conf")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	return viper.ReadInConfig()
}
