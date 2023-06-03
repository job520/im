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
		ctx.JSON(http.StatusOK, internal.Response{
			Ok:  false,
			Msg: err.Error(),
		})
		return
	}
	if err := service.Validate(param); err != nil {
		ctx.JSON(http.StatusOK, internal.Response{
			Ok:  false,
			Msg: err.Error(),
		})
		return
	}
	ok, err := service.Register(ctx, param.Username, param.Password)
	if err != nil {
		ctx.JSON(http.StatusOK, internal.Response{
			Ok:  false,
			Msg: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, internal.Response{
		Ok:  ok,
		Msg: "",
	})
	return
}

func Login(ctx *gin.Context) {
	var param internal.LoginRequest
	if err := ctx.ShouldBind(&param); err != nil {
		ctx.JSON(http.StatusOK, internal.Response{
			Ok:  false,
			Msg: err.Error(),
		})
		return
	}
	if err := service.Validate(param); err != nil {
		ctx.JSON(http.StatusOK, internal.Response{
			Ok:  false,
			Msg: err.Error(),
		})
		return
	}
	token, err := service.Login(ctx, param.Username, param.Password, param.Platform)
	if err != nil {
		ctx.JSON(http.StatusOK, internal.Response{
			Ok:  false,
			Msg: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, internal.Response{
		Ok:  true,
		Msg: "",
		Data: map[string]string{
			"token": token,
		},
	})
	return
}
