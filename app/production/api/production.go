package main

import (
	"erp/common/response"
	"erp/common/xcode"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/rest/httpx"

	"erp/app/production/api/internal/config"
	"erp/app/production/api/internal/handler"
	"erp/app/production/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/production-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	// 自定义错误处理方法
	httpx.SetErrorHandler(xcode.ErrHandler)

	// 自定义成功返回，仅对JSON有用
	httpx.SetOkHandler(response.OkHandler)
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
