package employeedetaillogic

import (
	"context"

	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelEmployeeDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelEmployeeDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelEmployeeDetailLogic {
	return &DelEmployeeDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelEmployeeDetailLogic) DelEmployeeDetail(in *pb.DelEmployeeDetailReq) (*pb.DelEmployeeDetailResp, error) {
	// 按工号删除员工扩展信息
	if err := l.svcCtx.EmployeeDetailModel.Delete(l.ctx, in.Id); err != nil {
		return nil, err
	}
	return &pb.DelEmployeeDetailResp{}, nil
}
