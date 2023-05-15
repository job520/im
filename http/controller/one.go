package controller

import (
	"github.com/gin-gonic/gin"
	"im/http/service"
	"net/http"
)

func GetOne(ctx *gin.Context) {
	name := ctx.Query("name")
	result, err := service.GetOne(ctx, name)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
	}
	ctx.String(http.StatusOK, result)
}
