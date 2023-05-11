package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetOne(context *gin.Context) {
	context.String(http.StatusOK, "success")
}
