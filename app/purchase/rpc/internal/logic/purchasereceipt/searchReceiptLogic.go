package purchasereceiptlogic

import (
"context"

	"erp/app/purchase/rpc/internal/svc"
	"erp/app/purchase/rpc/internal/types"
	"erp/app/purchase/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"erp/app/purchase/rpc/internal/code"
)

type SearchReceiptLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchReceiptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchReceiptLogic {
	return &SearchReceiptLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询入库单
func (l *SearchReceiptLogic) SearchReceipt(in *pb.SearchReceiptReq) (*pb.SearchReceiptResp, error) {
	params := &types.SearchReceiptParams{
		SearchComm:  types.SearchComm{Page: in.Page, Limit: in.Limit},
		ReceiptNo:   in.ReceiptNo,
		OrderId:     in.OrderId,
		WarehouseId: in.WarehouseId,
	}
	receipts, total, err := l.svcCtx.PurchaseReceiptModel.Search(l.ctx, params)
	if err != nil {

		return nil, code.GetReceiptFail

	}

	var receiptsWithDetails []*pb.PurchaseReceiptWithDetails
	for _, receipt := range receipts {
		receiptsWithDetail := &pb.PurchaseReceiptWithDetails{
			Receipts: &pb.PurchaseReceipt{
				Id:        receipt.Id,
				ReceiptNo: receipt.ReceiptNo,
				OrderId: func() int64 {
					if receipt.OrderId.Valid {
						return receipt.OrderId.Int64
					}
					return 0
				}(),
				WarehouseId:   receipt.WarehouseId,
				ReceiptDate:   receipt.ReceiptDate.Unix(),
				TotalQuantity: receipt.TotalQuantity,
				TotalAmount:   receipt.TotalAmount,
				Status:        receipt.Status,
				CreatedAt:     receipt.CreatedAt.Unix(),
				CreatedBy:     receipt.CreatedBy,
			},
			Total:   0,
			Details: nil,
		}

		details, err := l.svcCtx.PurchaseReceiptDetailModel.ListByReceiptId(l.ctx, receipt.Id)
		if err != nil {

			return nil, code.GetReceiptFail

		}
		receiptsWithDetail.Total = int64(len(details))
		for _, detail := range details {
			receiptsWithDetail.Details = append(receiptsWithDetail.Details, &pb.PurchaseReceiptDetail{
				Id:           detail.Id,
				ReceiptId:    detail.ReceiptId,
				ProductId:    detail.ProductId,
				ProductName:  detail.ProductName.String,
				CategoryType: detail.CategoryType,
				Quantity:     detail.Quantity,
				UnitPrice:    detail.UnitPrice,
				Amount:       detail.Amount,
				BatchId:      func() int64 { if detail.BatchId.Valid { return detail.BatchId.Int64 }; return 0 }(),
			})
		}
		receiptsWithDetails = append(receiptsWithDetails, receiptsWithDetail)
	}

	return &pb.SearchReceiptResp{
		ReceiptsWithDetails: receiptsWithDetails,
		Total:               total,
	}, nil
}