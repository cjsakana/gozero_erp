package purchaseReceipt

import (
	"context"
	"erp/app/purchase/api/internal/svc"
	"erp/app/purchase/api/internal/types"
	"erp/app/purchase/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateReceiptWithDetailsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateReceiptWithDetailsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateReceiptWithDetailsLogic {
	return &CreateReceiptWithDetailsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateReceiptWithDetailsLogic) CreateReceiptWithDetails(req *types.CreateReceiptWithDetailsReq) (resp *types.CreateReceiptWithDetailsResp, err error) {
	details := make([]*pb.ReceiptDetailInput, 0, len(req.Details))
	for _, d := range req.Details {
		productId, err := util.StringToInt64(d.ProductId)
		if err != nil {
			return nil, err
		}
		batchId, err := util.StringToInt64(d.BatchId)
		if err != nil {
			return nil, err
		}
		details = append(details, &pb.ReceiptDetailInput{
			ProductId:    productId,
			ProductName:  d.ProductName,
			CategoryType: d.CategoryType,
			Quantity:     d.Quantity,
			UnitPrice:    d.UnitPrice,
			Amount:       d.Amount,
			BatchId:      batchId,
		})
	}

	orderId, err := util.StringToInt64(req.OrderId)
	if err != nil {
		return nil, err
	}
	warehouseId, err := util.StringToInt64(req.WarehouseId)
	if err != nil {
		return nil, err
	}
	createdById, err := util.StringToInt64(req.CreatedById)
	if err != nil {
		return nil, err
	}

	ret, err := l.svcCtx.PurchaseRPC.CreateReceiptWithDetails(l.ctx, &pb.CreateReceiptWithDetailsReq{
		ReceiptNo:     req.ReceiptNo,
		OrderId:       orderId,
		WarehouseId:   warehouseId,
		ReceiptDate:   req.ReceiptDate,
		TotalQuantity: req.TotalQuantity,
		TotalAmount:   req.TotalAmount,
		Status:        req.Status,
		CreatedBy:     createdById,
		Details:       details,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.CreateReceiptWithDetailsResp{
		ReceiptId: util.Int64ToString(ret.ReceiptId),
	}
	return
}
