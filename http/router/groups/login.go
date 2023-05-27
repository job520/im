package groups

import (
	"github.com/gin-gonic/gin"
	"im/http/controller"
)

func LoadLogin(group *gin.RouterGroup) {
	group.PUT("/register", controller.Register) // 注册
	group.POST("/login", controller.Login)      // 登录
}
