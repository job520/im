package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// userid 和 Node 映射关系表
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

var rwlocker sync.RWMutex

type Node struct {
	Conn *websocket.Conn
	// 并行转串行,
	DataQueue chan []byte
	GroupSets set.Interface
}

// 定义命令行格式
const (
	CmdSingleMsg = 10
	CmdRoomMsg   = 11
	CmdHeart     = 0
)

type Message struct {
	Id      int64  `json:"id,omitempty" form:"id"`           // 消息ID
	Userid  int64  `json:"userid,omitempty" form:"userid"`   // 谁发的
	Cmd     int    `json:"cmd,omitempty" form:"cmd"`         // 群聊还是私聊
	Dstid   int64  `json:"dstid,omitempty" form:"dstid"`     // 对端用户ID/群ID
	Media   int    `json:"media,omitempty" form:"media"`     // 消息按照什么样式展示
	Content string `json:"content,omitempty" form:"content"` // 消息的内容
	Pic     string `json:"pic,omitempty" form:"pic"`         // 预览图片
	Url     string `json:"url,omitempty" form:"url"`         // 服务的URL
	Memo    string `json:"memo,omitempty" form:"memo"`       // 简单描述
	Amount  int    `json:"amount,omitempty" form:"amount"`   // 其他和数字相关的
}

// 后端调度逻辑处理
func dispatch(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		log.Println(err.Error())
		return
	}
	switch msg.Cmd {
	case CmdSingleMsg:
		sendMsg(msg.Dstid, data)
	case CmdRoomMsg:
		for _, v := range clientMap {
			if v.GroupSets.Has(msg.Dstid) {
				v.DataQueue <- data
			}
		}
	case CmdHeart:
		// 检测客户端的心跳
	}
}

// 发送逻辑
func sendproc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue: // 收到客户端消息
			err := node.Conn.WriteMessage(websocket.TextMessage, data) // 服务端回应客户端消息
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	}
}

// 接收逻辑
func recvproc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage() // 监听服务端接收到的消息
		if err != nil {
			log.Println(err.Error())
			return
		}

		dispatch(data) // 消息处理
		// todo: 对 data 进一步处理
		fmt.Printf("recv<=%s", data)
	}
}

// 发送消息,发送到消息的管道
func sendMsg(userId int64, msg []byte) {
	rwlocker.RLock()
	node, ok := clientMap[userId]
	rwlocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}

func checkToken(userId int64, token string) bool {
	return true
}

func Chat(ctx *gin.Context) {
	id := ctx.Query("id")
	token := ctx.Query("token")
	userId, _ := strconv.ParseInt(id, 10, 64)
	// 校验 token 是否合法
	islegal := checkToken(userId, token)

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
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}

	// 获取用户全部群 Id
	comIds := []int{1, 2, 3}
	for _, v := range comIds {
		node.GroupSets.Add(v)
	}

	rwlocker.Lock()
	clientMap[userId] = node
	rwlocker.Unlock()

	// 开启协程处理发送逻辑
	go sendproc(node)

	// 开启协程完成接收逻辑
	go recvproc(node)

	sendMsg(userId, []byte("welcome!"))
}

func main() {
	r := gin.Default()
	r.GET("/chat", Chat)
}

// todo: clientMap 拆分
// todo: 添加心跳处理逻辑
// todo: 添加异常关闭的逻辑（signal.Notify...）
// todo: 在图片上传的时候做 hash校验，如果资源文件已经存在了，直接将 url 返回给客户端
// todo: 拆分为分布式服务
