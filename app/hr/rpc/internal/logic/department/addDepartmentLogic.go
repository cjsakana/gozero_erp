package departmentlogic

import (
	"context"
	"database/sql"
	"erp/app/hr/rpc/internal/code"
	"erp/app/hr/rpc/internal/model"
	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/pb"
	"erp/common/util"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddDepartmentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddDepartmentLogic {
	return &AddDepartmentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------部门表-----------------------
func (l *AddDepartmentLogic) AddDepartment(in *pb.AddDepartmentReq) (*pb.AddDepartmentResp, error) {
	id := util.GenerateSnowflake()
	_, err := l.svcCtx.DepartmentModel.Insert(l.ctx, &model.Department{
		Id:          id,
		Name:        in.Name,
		ParentId:    sql.NullInt64{Int64: in.ParentId, Valid: true},
		Code:        sql.NullString{String: in.Code, Valid: true},
		ManagerId:   sql.NullInt64{Int64: in.ManagerId, Valid: true},
		ManagerNo:   sql.NullString{String: in.ManagerNo, Valid: true},
		ManagerName: in.ManagerName,
	})
	if err != nil {

		return nil, code.AddDepartmentFail

	}
	return &pb.AddDepartmentResp{
		Id: id,
	}, nil
}
