package departmentlogic

import (
	"context"
	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDepartmentByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDepartmentByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDepartmentByIdLogic {
	return &GetDepartmentByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDepartmentByIdLogic) GetDepartmentById(in *pb.GetDepartmentByIdReq) (*pb.GetDepartmentByIdResp, error) {
	one, err := l.svcCtx.DepartmentModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetDepartmentByIdResp{
		Department: &pb.Department{
			Id:          one.Id,
			Name:        one.Name,
			ParentId:    one.ParentId.Int64,
			Code:        one.Code.String,
			ManagerId:   one.ManagerId.Int64,
			ManagerNo:   one.ManagerNo.String,
			ManagerName: one.ManagerName,
			CreatedAt:   one.CreatedAt.Unix(),
			UpdatedAt:   one.UpdatedAt.Unix(),
		},
	}, nil
}
