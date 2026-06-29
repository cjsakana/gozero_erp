package receiptrecordlogic

import (
	"context"
	"database/sql"
	"erp/common/util"
	"time"

	"erp/app/finance/rpc/internal/code"
	"erp/app/finance/rpc/internal/model"
	"erp/app/finance/rpc/internal/svc"
	"erp/app/finance/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddReceiptRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	Logger logx.Logger
}

func NewAddReceiptRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddReceiptRecordLogic {
	return &AddReceiptRecordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------receiptRecord-----------------------
func (l *AddReceiptRecordLogic) AddReceiptRecord(in *pb.AddReceiptRecordReq) (*pb.AddReceiptRecordResp, error) {
	id := util.GenerateSnowflake()
	receiptNo := util.GenerateNo("RCP")
	data := &model.ReceiptRecord{
		Id:            id,
		ReceiptNo:     receiptNo,
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

	_, err := l.svcCtx.ReceiptRecordModel.Insert(l.ctx, data)
	if err != nil {
		return nil, code.AddReceiptRecordFail
	}

	return &pb.AddReceiptRecordResp{Id: id}, nil
}
