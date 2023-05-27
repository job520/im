package service

import "github.com/gin-gonic/gin"

func Register(ctx *gin.Context, name string) (string, error) {
	return name, nil
}

func Login(ctx *gin.Context, name string) (string, error) {
	return name, nil
}
