package purchaseorderlogic

import (
	"context"
	"erp/common/util"

	"erp/app/purchase/rpc/internal/svc"
	"erp/app/purchase/rpc/internal/types"
	"erp/app/purchase/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"erp/app/purchase/rpc/internal/code"
)

type CreateOrderWithDetailsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderWithDetailsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderWithDetailsLogic {
	return &CreateOrderWithDetailsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 直接创建采购订单及明细（事务）
func (l *CreateOrderWithDetailsLogic) CreateOrderWithDetails(in *pb.CreateOrderWithDetailsReq) (*pb.CreateOrderWithDetailsResp, error) {
	// 生成主表雪花ID
	orderId := util.GenerateSnowflake()

	param := &types.CreateOrderWithDetailsParam{
		OrderNo:      in.OrderNo,
		SupplierId:   in.SupplierId,
		OrderDate:    in.OrderDate,
		ExpectedDate: in.ExpectedDate,
		TotalAmount:  in.TotalAmount,
		Status:       in.Status,
		PurchaserId:  in.PurchaserId,
	}
	// 为每个明细生成雪花ID
	for _, d := range in.Details {
		param.Details = append(param.Details, types.OrderDetailParam{
			Id:           util.GenerateSnowflake(),
			ProductId:    d.ProductId,
			ProductName:  d.ProductName,
			CategoryType: d.CategoryType,
			Quantity:     d.Quantity,
			UnitPrice:    d.UnitPrice,
			Amount:       d.Amount,
			Remark:       d.Remark,
		})
	}

	err := l.svcCtx.PurchaseOrderModel.CreateWithDetails(l.ctx, orderId, param)
	if err != nil {

		return nil, code.CreateOrderFail

	}
	return &pb.CreateOrderWithDetailsResp{OrderId: orderId}, nil
}