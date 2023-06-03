package service

import (
	"github.com/gin-gonic/gin"
	"im/http/model"
)

func Register(ctx *gin.Context, username, password string) (bool, error) {
	return model.Register(ctx, username, password)
}

func Login(ctx *gin.Context, username, password string, platform int) (string, error) {
	return model.Login(ctx, username, password, platform)
}
