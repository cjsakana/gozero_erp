package payroll

import (
	"context"
	"erp/app/hr/rpc/pb"
	"erp/common/util"

	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubmitToFinanceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubmitToFinanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitToFinanceLogic {
	return &SubmitToFinanceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubmitToFinanceLogic) SubmitToFinance(req *types.SubmitToFinanceRequest) (resp *types.SubmitToFinanceResponse, err error) {
	var successCount, failCount int64
	var items []string
	for _, idStr := range req.Ids {
		id, err := util.StringToInt64(idStr)
		if err != nil {
			return nil, err
		}
		_, err = l.svcCtx.HrRPC.UpdatePayrollRecord(l.ctx, &pb.UpdatePayrollRecordReq{
			Id:     id,
			Status: 2,
		})
		if err != nil {
			failCount++
			items = append(items, idStr)
		} else {
			successCount++
		}
	}

	resp = &types.SubmitToFinanceResponse{
		SuccessCount: successCount,
		FailCount:    failCount,
		FailedIds:    items,
	}

	return
}
