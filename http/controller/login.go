package controller

import (
	"github.com/gin-gonic/gin"
	"im/http/controller/internal"
	"im/http/service"
	"net/http"
)

func Register(ctx *gin.Context) {
	var param internal.RegisterRequest
	if err := ctx.ShouldBind(&param); err != nil {
		ctx.JSON(http.StatusOK, internal.RegisterResponse{
			Ok:  false,
			Msg: err.Error(),
		})
		return
	}
	ok, err := service.Register(ctx, param.Username, param.Password)
	if err != nil {
		ctx.JSON(http.StatusOK, internal.RegisterResponse{
			Ok:  false,
			Msg: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, internal.RegisterResponse{
		Ok:  ok,
		Msg: "",
	})
	return
}

func Login(ctx *gin.Context) {
	name := ctx.Query("name")
	result, err := service.Login(ctx, name)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
	}
	ctx.String(http.StatusOK, result)
}
