package svc

import (
	"erp/app/hr/rpc/internal/config"
	"erp/app/hr/rpc/internal/model"
	"erp/app/user/rpc/user"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config                   config.Config
	AttendanceRecordModel    model.AttendanceRecordModel
	AttendanceReplenishModel model.AttendanceReplenishModel
	DepartmentModel          model.DepartmentModel
	EmployeeDetailModel      model.EmployeeDetailModel
	LeaveApplicationModel    model.LeaveApplicationModel
	PayrollRecordModel       model.PayrollRecordModel
	PositionModel            model.PositionModel
	ResignedApplicationModel model.ResignedApplicationModel
	BizRedis                 *redis.Redis
	UserRPC                  user.UserZrpcClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	rds, err := redis.NewRedis(redis.RedisConf{
		Host: c.BizRedis.Host,
		Pass: c.BizRedis.Pass,
		Type: c.BizRedis.Type,
	})
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:                   c,
		AttendanceRecordModel:    model.NewAttendanceRecordModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		AttendanceReplenishModel: model.NewAttendanceReplenishModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		DepartmentModel:          model.NewDepartmentModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		EmployeeDetailModel:      model.NewEmployeeDetailModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		LeaveApplicationModel:    model.NewLeaveApplicationModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		PayrollRecordModel:       model.NewPayrollRecordModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		PositionModel:            model.NewPositionModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		ResignedApplicationModel: model.NewResignedApplicationModel(sqlx.NewMysql(c.DataSource), c.CacheRedis),
		BizRedis:                 rds,
		UserRPC:                  user.NewUserZrpcClient(zrpc.MustNewClient(c.UserRPC)),
	}
}
