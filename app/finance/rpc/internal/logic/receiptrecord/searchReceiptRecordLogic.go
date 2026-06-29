package receiptrecordlogic

import (
	"context"

	"erp/app/finance/rpc/internal/model"
	"erp/app/finance/rpc/internal/svc"
	"erp/app/finance/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchReceiptRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	Logger logx.Logger
}

func NewSearchReceiptRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchReceiptRecordLogic {
	return &SearchReceiptRecordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchReceiptRecordLogic) SearchReceiptRecord(in *pb.SearchReceiptRecordReq) (*pb.SearchReceiptRecordResp, error) {
	list, total, err := l.svcCtx.ReceiptRecordModel.Search(l.ctx,
		in.ReceiptNo, in.ReceiptMethod, "", in.CustomerId, in.ReceiptType, in.Status, in.Page, in.Limit,
	)
	if err != nil {
		return nil, model.ErrNotFound
	}

	items := make([]*pb.ReceiptRecord, 0, len(list))
	for _, rr := range list {
		items = append(items, &pb.ReceiptRecord{
			Id:            rr.Id,
			ReceiptNo:     rr.ReceiptNo,
			CustomerId:    rr.CustomerId,
			OrderId:       rr.OrderId.Int64,
			ReceiptType:   rr.ReceiptType,
			Amount:        rr.Amount,
			ReceiptDate:   rr.ReceiptDate.Unix(),
			ReceiptMethod: rr.ReceiptMethod.String,
			Status:        rr.Status,
			VerifyStatus:  rr.VerifyStatus,
			OperatorId:    rr.OperatorId,
			CreatedAt:     rr.CreatedAt.Unix(),
		})
	}
	return &pb.SearchReceiptRecordResp{ReceiptRecord: items, Total: total}, nil
}
