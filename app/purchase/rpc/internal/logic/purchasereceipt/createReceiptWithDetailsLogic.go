package purchasereceiptlogic

import (
	"context"
	"erp/common/util"

	"erp/app/purchase/rpc/internal/svc"
	"erp/app/purchase/rpc/internal/types"
	"erp/app/purchase/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"erp/app/purchase/rpc/internal/code"
)

type CreateReceiptWithDetailsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateReceiptWithDetailsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateReceiptWithDetailsLogic {
	return &CreateReceiptWithDetailsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 直接创建入库单及明细（事务）
func (l *CreateReceiptWithDetailsLogic) CreateReceiptWithDetails(in *pb.CreateReceiptWithDetailsReq) (*pb.CreateReceiptWithDetailsResp, error) {
	// 生成主表雪花ID
	receiptId := util.GenerateSnowflake()

	param := &types.CreateReceiptWithDetailsParam{
		ReceiptNo:     in.ReceiptNo,
		OrderId:       in.OrderId,
		WarehouseId:   in.WarehouseId,
		ReceiptDate:   in.ReceiptDate,
		TotalQuantity: in.TotalQuantity,
		TotalAmount:   in.TotalAmount,
		Status:        in.Status,
		CreatedBy:     in.CreatedBy,
	}
	// 为每个明细生成雪花ID
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

	err := l.svcCtx.PurchaseReceiptModel.CreateWithDetails(l.ctx, receiptId, param)
	if err != nil {

		return nil, code.CreateReceiptFail

	}
	return &pb.CreateReceiptWithDetailsResp{ReceiptId: receiptId}, nil
}