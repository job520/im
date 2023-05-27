package router

import (
	"github.com/gin-gonic/gin"
	"im/http/middleware"
	"im/http/router/groups"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Test())
	groupLogin := r.Group("/login")
	groups.LoadLogin(groupLogin)
	return r
}
