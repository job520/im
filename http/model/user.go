package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"im/http/driver"
)

type User struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
}

func Register(ctx *gin.Context, username, password string) (bool, error) {
	collection, err := driver.NewMongoCollection("user")
	if err != nil {
		return false, err
	}
	// 查询是否存在记录
	count, err := collection.CountDocuments(ctx, bson.M{
		"username": username,
	})
	if err != nil {
		return false, err
	}
	if count > 0 {
		return false, fmt.Errorf("该用户已注册")
	}
	_, err = collection.InsertOne(ctx, User{
		Username: username,
		Password: password,
	})
	if err != nil {
		return false, err
	}
	return true, nil
}
