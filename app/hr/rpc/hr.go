package main

import (
	"context"
	"erp/app/hr/rpc/internal/config"
	payrollrecordlogic "erp/app/hr/rpc/internal/logic/payrollrecord"
	attendancerecordServer "erp/app/hr/rpc/internal/server/attendancerecord"
	attendancereplenishServer "erp/app/hr/rpc/internal/server/attendancereplenish"
	departmentServer "erp/app/hr/rpc/internal/server/department"
	employeedetailServer "erp/app/hr/rpc/internal/server/employeedetail"
	leaveapplicationServer "erp/app/hr/rpc/internal/server/leaveapplication"
	payrollrecordServer "erp/app/hr/rpc/internal/server/payrollrecord"
	positionServer "erp/app/hr/rpc/internal/server/position"
	resignedapplicationServer "erp/app/hr/rpc/internal/server/resignedapplication"
	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/pb"
	"erp/common/interceptors"
	"flag"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/hr.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterAttendanceRecordServer(grpcServer, attendancerecordServer.NewAttendanceRecordServer(ctx))
		pb.RegisterAttendanceReplenishServer(grpcServer, attendancereplenishServer.NewAttendanceReplenishServer(ctx))
		pb.RegisterDepartmentServer(grpcServer, departmentServer.NewDepartmentServer(ctx))
		pb.RegisterEmployeeDetailServer(grpcServer, employeedetailServer.NewEmployeeDetailServer(ctx))
		pb.RegisterLeaveApplicationServer(grpcServer, leaveapplicationServer.NewLeaveApplicationServer(ctx))
		pb.RegisterPayrollRecordServer(grpcServer, payrollrecordServer.NewPayrollRecordServer(ctx))
		pb.RegisterPositionServer(grpcServer, positionServer.NewPositionServer(ctx))
		pb.RegisterResignedApplicationServer(grpcServer, resignedapplicationServer.NewResignedApplicationServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	// 自定义拦截器
	s.AddUnaryInterceptors(interceptors.ServerErrorInterceptor())

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)

	// 启动定时任务：每月7号 02:00 生成上月工资
	if c.EnablePayrollCron {
		c := cron.New(cron.WithLocation(mustLoadLocation("Asia/Shanghai")))
		// cron表达式: 0 2 7 * * 表示每月7号凌晨2点执行
		_, err := c.AddFunc("0 2 7 * *", func() {
			l := payrollrecordlogic.NewGenerateMonthlyPayrollLogic(context.Background(), ctx)
			_, _ = l.GenerateMonthlyPayroll(&pb.GenerateMonthlyPayrollReq{})
		})
		if err != nil {
			fmt.Printf("Failed to add cron job: %v\n", err)
		} else {
			c.Start()
			fmt.Println("Payroll cron job started: runs at 02:00 on the 7th of every month (Asia/Shanghai)")
		}
	}

	s.Start()
}

// mustLoadLocation 加载时区，失败时panic
func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(fmt.Sprintf("Failed to load location %s: %v", name, err))
	}
	return loc
}
