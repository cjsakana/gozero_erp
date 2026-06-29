package employeeDetail

import (
	"context"
	"erp/app/hr/rpc/pb"
	"erp/common/util"

	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdjustSalaryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdjustSalaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdjustSalaryLogic {
	return &AdjustSalaryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdjustSalaryLogic) AdjustSalary(req *types.AdjustSalaryRequest) (resp *types.EmptyResponse, err error) {
	employeeId, err := util.StringToInt64(req.EmployeeId)
	_, err = l.svcCtx.HrRPC.EmployeeDetailZrpcClient.UpdateEmployeeDetail(l.ctx, &pb.UpdateEmployeeDetailReq{
		Id:     employeeId,
		Salary: req.NewSalary,
	})
	if err != nil {
		return nil, err
	}

	return
}
