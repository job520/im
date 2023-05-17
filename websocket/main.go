package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"im/websocket/global"
	"im/websocket/logic"
	"im/websocket/utils"
	"log"
	"net/http"
	"strconv"
)

func Chat(ctx *gin.Context) {
	id := ctx.Query("id")
	token := ctx.Query("token")
	userId, _ := strconv.Atoi(id)
	// 校验 token 是否合法
	islegal := logic.CheckToken(userId, token)

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
		DataQueue: make(chan string, 10),
		GroupSets: *utils.NewSet(),
	}

	// 获取用户全部群 Id
	groupIDArr := []int{1, 2, 3}
	for _, v := range groupIDArr {
		node.GroupSets.Add(v)
	}

	global.RwLocker.Lock()
	global.ClientMap[userId] = node
	global.RwLocker.Unlock()

	// 开启协程处理发送逻辑
	go logic.Send(node)

	// 开启协程完成接收逻辑
	go logic.Receive(node)

	logic.SendMsg(userId, "hello from server!")
}

func main() {
	r := gin.Default()
	r.GET("/chat", Chat)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

// todo: 业务逻辑拆分
// todo: 消息拆分 -- 单聊消息、群聊消息、系统消息
// todo: clientMap 拆分
// todo: 添加心跳处理逻辑
// todo: 添加异常关闭的逻辑（signal.Notify...）
// todo: 在图片上传的时候做 hash校验，如果资源文件已经存在了，直接将 url 返回给客户端
// todo: 拆分为分布式服务
// todo: 升级为 wss
