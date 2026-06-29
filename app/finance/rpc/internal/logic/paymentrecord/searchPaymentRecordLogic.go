package paymentrecordlogic

import (
	"context"

	"erp/app/finance/rpc/internal/svc"
	"erp/app/finance/rpc/pb"

	"erp/app/finance/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type SearchPaymentRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchPaymentRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchPaymentRecordLogic {
	return &SearchPaymentRecordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchPaymentRecordLogic) SearchPaymentRecord(in *pb.SearchPaymentRecordReq) (*pb.SearchPaymentRecordResp, error) {
	paymentRecords, total, err := l.svcCtx.PaymentRecordModel.Search(
		l.ctx,
		in.PaymentNo,
		in.PaymentMethod,
		in.SupplierId,
		in.PaymentType,
		in.Status,
		in.Page,
		in.Limit,
	)
	if err != nil {

		return nil, code.GetPaymentRecordFail

	}

	var pbPaymentRecords []*pb.PaymentRecord
	for _, pr := range paymentRecords {
		pbPaymentRecords = append(pbPaymentRecords, &pb.PaymentRecord{
			Id:            pr.Id,
			PaymentNo:     pr.PaymentNo,
			SupplierId:    pr.SupplierId,
			OrderId:       pr.OrderId.Int64,
			PaymentType:   pr.PaymentType,
			Amount:        pr.Amount,
			PaymentDate:   pr.PaymentDate.Unix(),
			PaymentMethod: pr.PaymentMethod.String,
			Status:        pr.Status,
			VerifyStatus:  pr.VerifyStatus,
			OperatorId:    pr.OperatorId,
			CreatedAt:     pr.CreatedAt.Unix(),
		})
	}

	return &pb.SearchPaymentRecordResp{
		PaymentRecord: pbPaymentRecords,
		Total:         total,
	}, nil
}
