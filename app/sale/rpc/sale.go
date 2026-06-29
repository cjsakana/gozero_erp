package main

import (
	"erp/common/interceptors"
	"flag"
	"fmt"

	"erp/app/sale/rpc/internal/config"
	salesdeliveryServer "erp/app/sale/rpc/internal/server/salesdelivery"
	salesorderServer "erp/app/sale/rpc/internal/server/salesorder"
	"erp/app/sale/rpc/internal/svc"
	"erp/app/sale/rpc/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/sale.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterSalesDeliveryServer(grpcServer, salesdeliveryServer.NewSalesDeliveryServer(ctx))
		pb.RegisterSalesOrderServer(grpcServer, salesorderServer.NewSalesOrderServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	// 自定义拦截器
	s.AddUnaryInterceptors(interceptors.ServerErrorInterceptor())

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
