package department

import (
	"context"
	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"
	"erp/app/hr/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDepartmentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDepartmentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDepartmentListLogic {
	return &GetDepartmentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDepartmentListLogic) GetDepartmentList(req *types.GetDepartmentListRequest) (resp *types.GetDepartmentListResponse, err error) {
	parentId, err := util.StringToInt64(req.ParentId)
	if err != nil {
		return nil, err
	}
	
	var queue []int64
	// 没有，即0，不会搜索的
	queue = append(queue, parentId)

	resp = &types.GetDepartmentListResponse{}

	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		ret, err := l.svcCtx.HrRPC.DepartmentZrpcClient.SearchDepartment(l.ctx, &pb.SearchDepartmentReq{
			Limit:    -1,
			Name:     req.Name,
			ParentId: id,
			Code:     req.Code,
		})
		if err != nil {
			return nil, err
		}
		for _, department := range ret.Department {
			resp.List = append(resp.List, &types.Department{
				Id:          util.Int64ToString(department.Id),
				Name:        department.Name,
				ParentId:    util.Int64ToString(department.ParentId),
				Code:        department.Code,
				ManagerId:   util.Int64ToString(department.ManagerId),
				ManagerNo:   department.ManagerNo,
				ManagerName: department.ManagerName,
				CreatedAt:   department.CreatedAt,
				UpdatedAt:   department.UpdatedAt,
			})
			queue = append(queue, department.Id)
		}
	}
	resp.Total = int64(len(resp.List))
	return
}
