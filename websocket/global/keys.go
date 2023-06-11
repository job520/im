package global

const (
	RedisJwtKey    = "jwt:%s:%d"    // jwt:userId:platform
	RedisStatusKey = "status:%s:%d" // status:userId:platform
	EtcdWsDir      = "/websocket/"
	EtcdRpcDir     = "/rpc/"
)
