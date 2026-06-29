package departmentlogic

import (
	"context"
	"erp/app/hr/rpc/internal/svc"
	types2 "erp/app/hr/rpc/internal/types"
	"erp/app/hr/rpc/pb"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type DelDepartmentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelDepartmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelDepartmentLogic {
	return &DelDepartmentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelDepartmentLogic) DelDepartment(in *pb.DelDepartmentReq) (*pb.DelDepartmentResp, error) {
	_, total, err := l.svcCtx.DepartmentModel.Search(l.ctx, &types2.SearchDepartmentParams{
		SearchCom: types2.SearchCom{
			Page:  0,
			Limit: -1,
		},
		ParentId: in.Id,
	})
	switch {
	case total > 0:
		return &pb.DelDepartmentResp{}, errors.New("删除失败，存在子部门")
	case err == nil:
		// 找不到即没有子部门了
		err = l.svcCtx.DepartmentModel.Delete(l.ctx, in.Id)
		if err != nil {
			return nil, err
		}
		return &pb.DelDepartmentResp{}, nil
	default:
		return &pb.DelDepartmentResp{}, errors.New("删除失败，未知错误")
	}
}
