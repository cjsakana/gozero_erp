package payroll

import (
	"context"
	"erp/app/hr/rpc/pb"
	"erp/common/util"
	"time"

	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExecutePaymentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExecutePaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExecutePaymentLogic {
	return &ExecutePaymentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExecutePaymentLogic) ExecutePayment(req *types.ExecutePaymentRequest) (resp *types.ExecutePaymentResponse, err error) {
	var successCount, failCount int64
	var items []string
	for _, idStr := range req.Ids {
		id, err := util.StringToInt64(idStr)
		if err != nil {
			return nil, err
		}
		_, err = l.svcCtx.HrRPC.UpdatePayrollRecord(l.ctx, &pb.UpdatePayrollRecordReq{
			Id:        id,
			PaymentAt: time.Now().Unix(),
			Status:    4,
		})
		if err != nil {
			failCount++
			items = append(items, idStr)
			// 应该是单独api
			_, err = l.svcCtx.HrRPC.UpdatePayrollRecord(l.ctx, &pb.UpdatePayrollRecordReq{
				Id:        id,
				PaymentAt: time.Now().Unix(),
				Status:    5,
			})
			if err != nil {
				l.Logger.Errorf("UpdatePayrollRecord failed: %v", err)
			}
		} else {
			successCount++
		}
	}

	resp = &types.ExecutePaymentResponse{
		SuccessCount: successCount,
		FailCount:    failCount,
		FailedIds:    items,
	}

	return
}
