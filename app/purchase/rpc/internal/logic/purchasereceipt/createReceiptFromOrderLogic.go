package purchasereceiptlogic

import (
	"context"
	"erp/common/util"

	"erp/app/purchase/rpc/internal/svc"
	"erp/app/purchase/rpc/internal/types"
	"erp/app/purchase/rpc/pb"

	"erp/app/purchase/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateReceiptFromOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateReceiptFromOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateReceiptFromOrderLogic {
	return &CreateReceiptFromOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 从采购订单创建入库单（通过 order_id 获取订单明细数据，保证一致性）
func (l *CreateReceiptFromOrderLogic) CreateReceiptFromOrder(in *pb.CreateReceiptFromOrderReq) (*pb.CreateReceiptFromOrderResp, error) {
	// 生成入库单雪花ID
	receiptId := util.GenerateSnowflake()

	// 获取订单明细
	orderDetails, err := l.svcCtx.PurchaseOrderDetailModel.ListByOrderId(l.ctx, in.OrderId)
	if err != nil {

		return nil, code.CreateOrderFail

	}

	param := &types.CreateReceiptFromOrderParam{
		OrderId:     in.OrderId,
		ReceiptNo:   in.ReceiptNo,
		WarehouseId: in.WarehouseId,
		ReceiptDate: in.ReceiptDate,
		CreatedBy:   in.CreatedBy,
	}

	// 如果提供了明细覆盖，则使用覆盖的明细；否则使用订单的未收货明细
	if len(in.Details) > 0 {
		for _, d := range in.Details {
			param.Details = append(param.Details, types.ReceiptDetailParam{
				Id:           util.GenerateSnowflake(),
				ProductId:    d.ProductId,
				ProductName:  d.ProductName,
				CategoryType: d.CategoryType,
				Quantity:     d.Quantity,
				UnitPrice:    d.UnitPrice,
				Amount:       d.Amount,
				BatchId:      d.BatchId,
			})
		}
	} else {
		// 使用订单的未收货明细（数量 = 订单数量 - 已收货数量）
		for _, d := range orderDetails {
			unreceivedQty := d.Quantity - d.ReceivedQty
			if unreceivedQty > 0 {
				param.Details = append(param.Details, types.ReceiptDetailParam{
					Id:           util.GenerateSnowflake(),
					ProductId:    d.ProductId.Int64,
					ProductName:  d.ProductName.String,
					CategoryType: d.CategoryType,
					Quantity:     unreceivedQty,
					UnitPrice:    d.UnitPrice,
					Amount:       unreceivedQty * d.UnitPrice,
					BatchId:      0, // 从订单创建时没有批次
				})
			}
		}
	}

	err = l.svcCtx.PurchaseReceiptModel.CreateFromOrder(l.ctx, receiptId, param)
	if err != nil {

		return nil, code.CreateOrderFail

	}

	// 更新订单明细的已收货数量
	for _, detail := range param.Details {
		orderDetail, err := l.svcCtx.PurchaseOrderDetailModel.ListByOrderId(l.ctx, in.OrderId)
		if err != nil {
			continue
		}
		for _, od := range orderDetail {
			if od.ProductId.Int64 == detail.ProductId {
				newReceivedQty := od.ReceivedQty + detail.Quantity
				_ = l.svcCtx.PurchaseOrderDetailModel.UpdateReceivedQty(l.ctx, in.OrderId, detail.ProductId, newReceivedQty)
				break
			}
		}
	}

	return &pb.CreateReceiptFromOrderResp{ReceiptId: receiptId}, nil
}
