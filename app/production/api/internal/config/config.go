package config

import (
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	Etcd           discov.EtcdConf
	AuthRPC        zrpc.RpcClientConf
	ProductionRPC  zrpc.RpcClientConf
	HrRPC          zrpc.RpcClientConf
	InventoryRPC   zrpc.RpcClientConf
	ProductRPC     zrpc.RpcClientConf
	BizRedis       redis.RedisConf
}
