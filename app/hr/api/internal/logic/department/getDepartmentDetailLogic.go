package department

import (
	"context"
	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"
	"erp/app/hr/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDepartmentDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDepartmentDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDepartmentDetailLogic {
	return &GetDepartmentDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDepartmentDetailLogic) GetDepartmentDetail(req *types.GetDepartmentDetailRequest) (resp *types.GetDepartmentDetailResponse, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	
	ret, err := l.svcCtx.HrRPC.DepartmentZrpcClient.GetDepartmentById(l.ctx, &pb.GetDepartmentByIdReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.GetDepartmentDetailResponse{
		Department: types.Department{
			Id:          util.Int64ToString(ret.Department.Id),
			Name:        ret.Department.Name,
			ParentId:    util.Int64ToString(ret.Department.ParentId),
			Code:        ret.Department.Code,
			ManagerId:   util.Int64ToString(ret.Department.ManagerId),
			ManagerNo:   ret.Department.ManagerNo,
			ManagerName: ret.Department.ManagerName,
			CreatedAt:   ret.Department.CreatedAt,
			UpdatedAt:   ret.Department.UpdatedAt,
		},
	}

	return
}
