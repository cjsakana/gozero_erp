package main

import (
	"erp/common/interceptors"
	"flag"
	"fmt"

	"erp/app/purchase/rpc/internal/config"
	purchaseorderServer "erp/app/purchase/rpc/internal/server/purchaseorder"
	purchasereceiptServer "erp/app/purchase/rpc/internal/server/purchasereceipt"
	purchaserequisitionServer "erp/app/purchase/rpc/internal/server/purchaserequisition"
	"erp/app/purchase/rpc/internal/svc"
	"erp/app/purchase/rpc/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/purchase.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterPurchaseRequisitionServer(grpcServer, purchaserequisitionServer.NewPurchaseRequisitionServer(ctx))
		pb.RegisterPurchaseOrderServer(grpcServer, purchaseorderServer.NewPurchaseOrderServer(ctx))
		pb.RegisterPurchaseReceiptServer(grpcServer, purchasereceiptServer.NewPurchaseReceiptServer(ctx))

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
