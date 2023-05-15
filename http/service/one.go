package service

import "github.com/gin-gonic/gin"

func GetOne(ctx *gin.Context, name string) (string, error) {
	return name, nil
}
