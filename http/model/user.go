package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"im/http/driver"
)

type User struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
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

func Login(ctx *gin.Context, username, password string, platform int) (string, error) {
	collection, err := driver.NewMongoCollection("user")
	if err != nil {
		return "", err
	}
	// 查询一条记录
	userInfoResult := collection.FindOne(ctx, bson.M{
		"username": username,
		"password": password,
	})
	if err := userInfoResult.Err(); err != nil {
		return "", err
	}
	var userInfo User
	if err := userInfoResult.Decode(&userInfo); err != nil {
		return "", err
	}
	return "", nil

	//// 查询是否存在记录
	//count, err := collection.CountDocuments(ctx, bson.M{
	//	"username": username,
	//	"password": password,
	//})
	//if err != nil {
	//	return "", err
	//}
	//if count == 0 {
	//	return "", fmt.Errorf("用户名或密码错误")
	//}
	//_, err = collection.InsertOne(ctx, User{
	//	Username: username,
	//	Password: password,
	//})
	//if err != nil {
	//	return false, err
	//}
	//return true, nil
}
