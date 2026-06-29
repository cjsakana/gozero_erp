package config

import (
	"erp/common/upload"
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
	Etcd         discov.EtcdConf
	AuthRPC      zrpc.RpcClientConf
	HrRPC        zrpc.RpcClientConf
	PurchaseRPC  zrpc.RpcClientConf
	SupplierRPC  zrpc.RpcClientConf
	ProductRPC   zrpc.RpcClientConf
	InventoryRPC zrpc.RpcClientConf
	BizRedis     redis.RedisConf
	R2Conf       upload.R2Conf
}
