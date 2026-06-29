package departmentlogic

import (
	"context"
	"erp/app/hr/rpc/internal/types"

	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchDepartmentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchDepartmentLogic {
	return &SearchDepartmentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchDepartmentLogic) SearchDepartment(in *pb.SearchDepartmentReq) (*pb.SearchDepartmentResp, error) {
	departments, total, err := l.svcCtx.DepartmentModel.Search(l.ctx, &types.SearchDepartmentParams{
		SearchCom: types.SearchCom{
			Page:  in.Page,
			Limit: in.Limit,
		},
		Name:     in.Name,
		ParentId: in.ParentId,
		Code:     in.Code,
	})
	if err != nil {
		return nil, err
	}

	var pbDepartments []*pb.Department
	for _, department := range departments {
		pbDepartments = append(pbDepartments, &pb.Department{
			Id:          department.Id,
			Name:        department.Name,
			ParentId:    department.ParentId.Int64,
			Code:        department.Code.String,
			ManagerId:   department.ManagerId.Int64,
			ManagerNo:   department.ManagerNo.String,
			ManagerName: department.ManagerName,
			CreatedAt:   department.CreatedAt.Unix(),
			UpdatedAt:   department.UpdatedAt.Unix(),
		})
	}
	return &pb.SearchDepartmentResp{
		Total:      total,
		Department: pbDepartments,
	}, nil
}
