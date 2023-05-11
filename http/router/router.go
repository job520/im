package router

import (
	"github.com/gin-gonic/gin"
	"im/http/router/groups"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	groupOne := r.Group("/one")
	groups.LoadOne(groupOne)
	return r
}
