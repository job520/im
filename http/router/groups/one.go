package groups

import (
	"github.com/gin-gonic/gin"
	"im/http/controller"
)

func LoadOne(group *gin.RouterGroup) {
	group.GET("/get", controller.GetOne)
}
