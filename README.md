# im

#### 介绍
im

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
8. redis: 缓存
9. [zinx]: tcp长链接服务
```