package driver

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"im/tcp/config"
	"time"
)

func NewMongoCollection(collectionName string) (*mongo.Collection, error) {
	username := config.Config.Mongodb.Username
	password := config.Config.Mongodb.Password
	host := config.Config.Mongodb.Host
	port := config.Config.Mongodb.Port
	database := config.Config.Mongodb.Database
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d", username, password, host, port)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri).SetConnectTimeout(5*time.Second))
	if err != nil {
		return nil, err
	}
	// 选择数据库
	db := client.Database(database)
	collection := db.Collection(collectionName)
	return collection, nil
}
