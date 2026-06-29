package main

import (
	"erp/common/interceptors"
	"flag"
	"fmt"

	"erp/app/product/rpc/internal/config"
	productServer "erp/app/product/rpc/internal/server/product"
	productcategoryServer "erp/app/product/rpc/internal/server/productcategory"
	productbatchServer "erp/app/product/rpc/internal/server/productbatch"
	"erp/app/product/rpc/internal/svc"
	"erp/app/product/rpc/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/product.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterProductServer(grpcServer, productServer.NewProductServer(ctx))
		pb.RegisterProductBatchServer(grpcServer, productbatchServer.NewProductBatchServer(ctx))
		pb.RegisterProductCategoryServer(grpcServer, productcategoryServer.NewProductCategoryServer(ctx))

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
