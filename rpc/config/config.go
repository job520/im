package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Config config

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatal("Fatal error config file:", err)
	}
	if err := viper.Unmarshal(&Config); err != nil {
		logrus.Fatal(err)
	}
}

type config struct {
	Server server `mapstructure:"server"`
}

type server struct {
	Address string `mapstructure:"address"`
}
