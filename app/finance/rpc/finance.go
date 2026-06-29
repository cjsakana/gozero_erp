package main

import (
	"flag"
	"fmt"

	"erp/app/finance/rpc/internal/config"
	fixedassetServer "erp/app/finance/rpc/internal/server/fixedasset"
	paymentrecordServer "erp/app/finance/rpc/internal/server/paymentrecord"
	receiptrecordServer "erp/app/finance/rpc/internal/server/receiptrecord"
	salarypaymentServer "erp/app/finance/rpc/internal/server/salarypayment"
	"erp/app/finance/rpc/internal/svc"
	"erp/app/finance/rpc/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/finance.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterFixedAssetServer(grpcServer, fixedassetServer.NewFixedAssetServer(ctx))
		pb.RegisterPaymentRecordServer(grpcServer, paymentrecordServer.NewPaymentRecordServer(ctx))
		pb.RegisterReceiptRecordServer(grpcServer, receiptrecordServer.NewReceiptRecordServer(ctx))
		pb.RegisterSalaryPaymentServer(grpcServer, salarypaymentServer.NewSalaryPaymentServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
