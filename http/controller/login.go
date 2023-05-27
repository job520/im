package controller

import (
	"github.com/gin-gonic/gin"
	"im/http/service"
	"net/http"
)

func Register(ctx *gin.Context) {
	name := ctx.Query("name")
	result, err := service.Register(ctx, name)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
	}
	ctx.String(http.StatusOK, result)
}

func Login(ctx *gin.Context) {
	name := ctx.Query("name")
	result, err := service.Login(ctx, name)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
	}
	ctx.String(http.StatusOK, result)
}
