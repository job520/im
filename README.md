# im

#### 介绍
im

#### 参考文档
[适合新手：从零开发一个IM服务端（基于Netty，有完整源码）](https://zhuanlan.zhihu.com/p/85758575)  

[GlideIM-Golang实现的高性能的分布式IM
](https://dengzii.com/go/5467dc73.html)  

[跟着源码学IM(六)：手把手教你用Go快速搭建高性能、可扩展的IM系统](http://www.52im.net/thread-2988-1-1.html)


#### 整体架构
| 角色                                           | 机器数量 | 功能                  |
|:---------------------------------------------|:-----|:--------------------|
| http 服务器                                     | 多台   | 处理登录、文件上传等业务        |
| websocket 服务器 + rpc 客户端（与 rpc服务端增加心跳检测与服务发现） | 多台   | 处理聊天业务              |
| rpc 服务器（与 rpc客户端增加心跳检测与服务发现）                 | 1台   | 处理消息转发业务            |
| redis 状态管理服务器                                | -    | 管理用户 & rpc 服务器的在线状态 |
| mongodb                                      | -    | 管理离线消息              |
| rabbitmq 消息队列服务器                             | -    | 防止消息重复              |
| etcd 服务发现与负载均衡服务器                            | -    | 服务发现&负载均衡           |


#### 项目计划

```text
### 第一步：搭建 http 服务器，处理注册、登录、文件上传等业务（分布式，跟 websocket 服务分开搭建）

### 第二步：搭建 websocket 服务，处理聊天业务（分布式，跟 http 服务分开搭建）

### 第三步：搭建 rpc 服务，处理 websocket 服务器之间的转发业务

### 第四步：细节处理
1. 不丢消息
2. 不重复
3. 不乱序
4. 消息加密
5. 离线消息存储
6. 防止离线消息重复推送
7. 图片上传做 hash 校验，如果已上传过，则直接返回之前上传的 url 地址
8. 接入 zipkin（链路追踪）
9. 优雅的停止服务

### 第五步：搭建 tcp 服务，使用 rpc 服务器与 websocket 进行业务整合与分发
```

#### 项目架构

```text
1. http 分布式服务
2. websocket 分布式服务
3. rpc 分布式服务，处理 websocket 服务器之间的转发业务
4. tcp 分布式服务
```

#### 用到的框架

```text
1. viper: 读取配置文件
2. gorilla/websocket: websocket长链接服务
3. rabbitmq: 消息队列，消息去重，消息顺序保证
4. 双向grpc + protobuf: 服务间通信
5. etcd: 服务发现，负载均衡
6. gin: 边沿业务，如登录、文件上传等
7. mongodb: 消息持久化存储
8. redis: 状态管理
9. [tcp socket]: tcp长链接服务
```

#### 项目依赖

```text
1. redis: 状态管理
2. mongodb: 消息持久化存储
3. etcd: 服务发现，负载均衡
4. rabbitmq: 消息队列，消息去重，消息顺序保证
```

#### TODO

```text
O 搭建开发环境：redis、mongodb、rabbitmq、etcd
O im/http: 注册/登录
O im/websocket: 启动时注册 server 到 etcd 中，并定时更新 TTL
O im/websocket: 建立连接后管理 【 userId:platform -> wsServer 】 在线状态（redis，添加 userId 与 connectorId(可以使用 ip + 端口) 的映射，需要与 websocket 客户端添加 TTL-心跳检测）
O im/rpc: 启动时注册 server 到 etcd 中，并定时更新 TTL
O im/websocket: 启动时连接到 rpc server（transfer 服务器）
O im/rpc: 记录与 rpc客户端（指 websocket 服务器）之间的连接句柄（map[ip:port]conn）
O im/http: 网关服务，获取 websocket 连接地址（从 etcd 中获取存活的 websocket server）
O im/websocket: 消息模型数据结构简化
O im/rpc: 仿照 im/websocket 服务器处理消息逻辑
O im/websocket: 消息转发（websocket服务器 -> rpc服务器）
O im/rpc: 消息转发（rpc服务器 -> websocket服务器）
O im/websocket: 简化消息模型（只保留单聊消息，去掉心跳检测和群聊消息）
O im/rpc: 简化消息模型（只保留转发消息，去掉心跳检测）
O im/tcp: 接入 tcp socket（仿照 websocket 项目，先做简单的单聊）
O 思考怎么接入群聊功能
X 细节处理：
    1. 不丢消息（ACK 机制）
    2. 不重复（消息添加全局唯一递增 ID）
    3. 不乱序（1.消息添加全局唯一递增 ID，客户端记录 lastID；2.上线时先推送离线消息，再将状态改为 online）
    4. 消息安全性（user_relations 表保存 AES 密钥）
    5. 离线消息存储（msg_type:消息类型(chat/ack)、to_user_id:接收消息的用户ID、has_read:是否已读）
    6. 防止离线消息重复推送（UPDATE `im_offline` SET has_read = true WHERE id = ${msg_id} AND has_read = false）
    7. 使用 etcd 制作配置中心
    8. 接入 zipkin（链路追踪）
X 使 tcp 服务器和 websocket 服务器之间可以实现消息实时转发
X 画出项目设计图
```