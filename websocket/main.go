package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"im/websocket/config"
	"im/websocket/global"
	"im/websocket/logic"
	"im/websocket/service"
	"im/websocket/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Chat(ctx *gin.Context) {
	token := ctx.GetHeader("token")
	// 校验 token 是否合法
	userId, platform, islegal := logic.CheckToken(token)

	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return islegal
		},
	}).Upgrade(ctx.Writer, ctx.Request, nil)

	if err != nil {
		log.Println(err.Error())
		return
	}
	// 获得 websocket 链接 conn
	node := &global.Node{
		Conn:      conn,
		MsgChan:   make(chan string, 10),
		GroupSets: *utils.NewSet(),
	}

	// 获取用户全部群 ID
	groupIDArr := []string{"group1", "group2", "group3"}
	for _, v := range groupIDArr {
		node.GroupSets.Add(v)
	}

	global.RwLocker.Lock()
	mapKey := fmt.Sprintf("%s:%d", userId, platform)
	global.ClientMap[mapKey] = node
	global.RwLocker.Unlock()

	// 开启协程处理发送逻辑
	go logic.Send(node)

	// 开启协程完成接收逻辑
	go logic.Receive(userId, platform, node)

	service.ReceiveSingleMsg(userId, platform, "hello from server!")
}

func main() {
	quit := make(chan os.Signal, 1) // 退出信号
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	r := gin.Default()
	r.GET("/chat", Chat)
	srv := &http.Server{
		Addr:    config.Config.Server.Address,
		Handler: r,
	}
	go shutdown(quit, srv)
	logic.RegisterWsServer() // 服务注册
	logic.RpcClient()        // 连接到 rpc 服务器（transfer 服务器）
	logrus.Infof("server running at:%s \n", config.Config.Server.Address)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func shutdown(quit chan os.Signal, srv *http.Server) {
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	logrus.Info("server shutdown...")
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal(err)
	}
	os.Exit(0)
}

// todo: 业务逻辑拆分
// todo: 添加 rpc客户端，处理消息转发业务
// todo: 消息拆分 -- 单聊消息、群聊消息、系统消息
// todo: clientMap 拆分
// todo: 添加心跳处理逻辑
// todo: 添加异常关闭的逻辑（signal.Notify...）
// todo: 在图片上传的时候做 hash校验，如果资源文件已经存在了，直接将 url 返回给客户端
// todo: 拆分为分布式服务
// todo: 升级为 wss
