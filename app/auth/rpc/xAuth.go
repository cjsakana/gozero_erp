package main

import (
	"erp/common/interceptors"
	"flag"
	"fmt"

	"erp/app/auth/rpc/internal/config"
	permissionServer "erp/app/auth/rpc/internal/server/permission"
	roleServer "erp/app/auth/rpc/internal/server/role"
	rolepermissionServer "erp/app/auth/rpc/internal/server/rolepermission"
	userroleServer "erp/app/auth/rpc/internal/server/userrole"
	"erp/app/auth/rpc/internal/svc"
	"erp/app/auth/rpc/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/xAuth.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 假如使用配置中心，与 conf.MustLoad(*configFile, &c) 互斥
	//c = config.PullConfig()

	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterPermissionServer(grpcServer, permissionServer.NewPermissionServer(ctx))
		pb.RegisterRoleServer(grpcServer, roleServer.NewRoleServer(ctx))
		pb.RegisterRolePermissionServer(grpcServer, rolepermissionServer.NewRolePermissionServer(ctx))
		pb.RegisterUserRoleServer(grpcServer, userroleServer.NewUserRoleServer(ctx))

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
