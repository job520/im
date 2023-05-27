package router

import (
	"github.com/gin-gonic/gin"
	"im/http/router/groups"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	groupLogin := r.Group("/login")
	groups.LoadLogin(groupLogin)
	return r
}
