package receiptrecordlogic

import (
	"context"
	"database/sql"
	"time"

	"erp/app/finance/rpc/internal/code"
	"erp/app/finance/rpc/internal/model"
	"erp/app/finance/rpc/internal/svc"
	"erp/app/finance/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateReceiptRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	Logger logx.Logger
}

func NewUpdateReceiptRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateReceiptRecordLogic {
	return &UpdateReceiptRecordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateReceiptRecordLogic) UpdateReceiptRecord(in *pb.UpdateReceiptRecordReq) (*pb.UpdateReceiptRecordResp, error) {
	data := &model.ReceiptRecord{
		Id:            in.Id,
		CustomerId:    in.CustomerId,
		OrderId:       sql.NullInt64{Int64: in.OrderId, Valid: in.OrderId != 0},
		ReceiptType:   in.ReceiptType,
		Amount:        in.Amount,
		ReceiptDate:   time.Unix(in.ReceiptDate, 0),
		ReceiptMethod: sql.NullString{String: in.ReceiptMethod, Valid: in.ReceiptMethod != ""},
		Status:        in.Status,
		VerifyStatus:  in.VerifyStatus,
		OperatorId:    in.OperatorId,
	}

	if err := l.svcCtx.ReceiptRecordModel.Update(l.ctx, data); err != nil {
		return nil, code.UpdateReceiptRecordFail
	}
	return &pb.UpdateReceiptRecordResp{}, nil
}
