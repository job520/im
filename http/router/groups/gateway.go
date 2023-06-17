package groups

import (
	"github.com/gin-gonic/gin"
	"im/http/controller"
)

func LoadGateWay(group *gin.RouterGroup) {
	group.GET("/get", controller.GetGateWay) // 获取长链接服务器地址
}
