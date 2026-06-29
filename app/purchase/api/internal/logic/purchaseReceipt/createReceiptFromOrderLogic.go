package purchaseReceipt

import (
	"context"
	"erp/app/purchase/api/internal/svc"
	"erp/app/purchase/api/internal/types"
	"erp/app/purchase/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateReceiptFromOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateReceiptFromOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateReceiptFromOrderLogic {
	return &CreateReceiptFromOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateReceiptFromOrderLogic) CreateReceiptFromOrder(req *types.CreateReceiptFromOrderReq) (resp *types.CreateReceiptFromOrderResp, err error) {
	// 从采购订单创建入库单：如果未提供明细，系统会自动从订单中获取未收货的明细数据
	// 如果提供了明细，则使用提供的明细（用于部分入库、调整数量等）
	var details []*pb.ReceiptDetailInput
	if req.Details != nil && len(req.Details) > 0 {
		details = make([]*pb.ReceiptDetailInput, 0, len(req.Details))
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
	}
	// 如果 details 为 nil 或空，RPC 层会从 orderId 对应的订单中自动获取未收货明细

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

	ret, err := l.svcCtx.PurchaseRPC.CreateReceiptFromOrder(l.ctx, &pb.CreateReceiptFromOrderReq{
		OrderId:     orderId,
		ReceiptNo:   req.ReceiptNo,
		WarehouseId: warehouseId,
		ReceiptDate: req.ReceiptDate,
		CreatedBy:   createdById,
		Details:     details, // nil 或空数组时，RPC 会从订单中获取未收货明细
	})
	if err != nil {
		return nil, err
	}

	resp = &types.CreateReceiptFromOrderResp{
		ReceiptId: util.Int64ToString(ret.ReceiptId),
	}
	return
}
