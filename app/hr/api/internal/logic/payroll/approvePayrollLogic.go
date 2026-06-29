package payroll

import (
	"context"
	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"
	"erp/app/hr/rpc/pb"
	"erp/common/util"
	"erp/common/xtypes"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApprovePayrollLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApprovePayrollLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApprovePayrollLogic {
	return &ApprovePayrollLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApprovePayrollLogic) ApprovePayroll(req *types.ApprovePayrollRecordRequest) (resp *types.ApprovePayrollRecordResponse, err error) {
	employeeId, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)
	if err != nil {
		return nil, err
	}

	var successCount, failCount int64
	var items []string
	for _, idStr := range req.Ids {
		id, err := util.StringToInt64(idStr)
		if err != nil {
			return nil, err
		}
		_, err = l.svcCtx.HrRPC.UpdatePayrollRecord(l.ctx, &pb.UpdatePayrollRecordReq{
			Id:           id,
			CalculatedBy: employeeId,
			CalculatedAt: time.Now().Unix(),
			Status:       1,
		})
		if err != nil {
			failCount++
			items = append(items, idStr)
		} else {
			successCount++
		}
	}

	resp = &types.ApprovePayrollRecordResponse{
		SuccessCount: successCount,
		FailCount:    failCount,
		FailedIds:    items,
	}
	return
}
