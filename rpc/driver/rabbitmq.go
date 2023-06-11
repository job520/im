package driver

import (
	"fmt"
	"github.com/streadway/amqp"
	"im/rpc/config"
)

func NewRabbitmqConnection() (*amqp.Connection, error) {
	username := config.Config.Rabbitmq.Username
	password := config.Config.Rabbitmq.Password
	host := config.Config.Rabbitmq.Host
	port := config.Config.Rabbitmq.Port
	url := fmt.Sprintf("amqp://%s:%s@%s:%d", username, password, host, port)
	connection, err := amqp.Dial(url)
	return connection, err
}
