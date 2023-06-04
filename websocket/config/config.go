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
	Server    server    `mapstructure:"server"`
	RpcServer rpcServer `mapstructure:"rpcServer"`
	Mongodb   mongodb   `mapstructure:"mongodb"`
	Redis     redis     `mapstructure:"redis"`
	Rabbitmq  rabbitmq  `mapstructure:"rabbitmq"`
	Etcd      etcd      `mapstructure:"etcd"`
	Jwt       jwt       `mapstructure:"jwt"`
}

type server struct {
	Address string `mapstructure:"address"`
}

type rpcServer struct {
	Address string `mapstructure:"address"`
}

type mongodb struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Username string `mapstructure:"username"`
	Database string `mapstructure:"database"`
}

type redis struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
}

type rabbitmq struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Username string `mapstructure:"username"`
}

type etcd struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Username string `mapstructure:"username"`
}

type jwt struct {
	EncryptKey  string `mapstructure:"encryptKey"`
	ExpireHours int    `mapstructure:"expireHours"`
}
