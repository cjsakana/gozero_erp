package departmentlogic

import (
	"context"
	"database/sql"
	"erp/app/hr/rpc/internal/model"
	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/pb"

	"erp/app/hr/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateDepartmentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDepartmentLogic {
	return &UpdateDepartmentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateDepartmentLogic) UpdateDepartment(in *pb.UpdateDepartmentReq) (*pb.UpdateDepartmentResp, error) {
	err := l.svcCtx.DepartmentModel.XUpdate(l.ctx, &model.Department{
		Id:          in.Id,
		Name:        in.Name,
		ParentId:    sql.NullInt64{Int64: in.ParentId, Valid: true},
		Code:        sql.NullString{String: in.Code, Valid: true},
		ManagerId:   sql.NullInt64{Int64: in.ManagerId, Valid: true},
		ManagerNo:   sql.NullString{String: in.ManagerNo, Valid: true},
		ManagerName: in.ManagerNo,
	})
	if err != nil {
		return nil, code.UpdateDepartmentFail
	}

	return &pb.UpdateDepartmentResp{}, nil
}
