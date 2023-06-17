package controller

import (
	"github.com/gin-gonic/gin"
	"im/http/controller/internal"
	"im/http/service"
	"net/http"
)

func GetGateWay(ctx *gin.Context) {
	serverAddress, err := service.GetGateWay(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, internal.Response{
			Ok:  false,
			Msg: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, internal.Response{
		Ok:   true,
		Msg:  "",
		Data: serverAddress,
	})
	return
}
