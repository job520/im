package router

import (
	"github.com/gin-gonic/gin"
	"im/http/middleware"
	"im/http/router/groups"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	groupLogin := r.Group("/login")
	groups.LoadLogin(groupLogin)
	gatewayGroup := r.Group("/gateway")
	gatewayGroup.Use(middleware.Test())
	groups.LoadGateWay(gatewayGroup)
	return r
}
