# im

#### 介绍
im

#### 参考文档
[适合新手：从零开发一个IM服务端（基于Netty，有完整源码）](https://zhuanlan.zhihu.com/p/85758575)  

[GlideIM-Golang实现的高性能的分布式IM
](https://dengzii.com/go/5467dc73.html)  

[跟着源码学IM(六)：手把手教你用Go快速搭建高性能、可扩展的IM系统](http://www.52im.net/thread-2988-1-1.html)


#### 整体架构
| 角色                                           | 功能                  |
|:---------------------------------------------|:--------------------|
| http 服务器                                     | 处理登录、文件上传等业务        |
| websocket 服务器 + rpc 客户端（与 rpc服务端增加心跳检测与服务发现） | 处理聊天业务              |
| rpc 服务器（与 rpc客户端增加心跳检测与服务发现）                 | 处理消息转发业务            |
| redis 状态管理服务器                                | 管理用户 & rpc 服务器的在线状态 |
| mongodb                                      | 管理离线消息              |
| rabbitmq 消息队列服务器                             | 防止消息重复              |


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
8. 接入 jaeger（链路追踪）
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
4. rpcx + protobuf: 服务间通信
5. etcd: 服务发现，负载均衡
6. gin: 边沿业务，如登录、文件上传等
7. mongodb: 消息持久化存储
8. redis: 状态管理
9. [zinx]: tcp长链接服务
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
im/http: 登录/注册
im/websocket: 登录后管理 websocket 状态（更新 redis，添加 userId 与 connectorId 的映射）
im/websocket: 管理 rpc 状态（更新 redis，添加 client 与 server 的映射）
im/rpc: 注册 server 到 etcd 中，并添加心跳检测
```