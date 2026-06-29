package purchaseOrder

import (
	"context"
	"erp/app/purchase/api/internal/svc"
	"erp/app/purchase/api/internal/types"
	"erp/app/purchase/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderWithDetailsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateOrderWithDetailsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderWithDetailsLogic {
	return &CreateOrderWithDetailsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateOrderWithDetailsLogic) CreateOrderWithDetails(req *types.CreateOrderWithDetailsReq) (resp *types.CreateOrderWithDetailsResp, err error) {
	details := make([]*pb.OrderDetailInput, 0, len(req.Details))
	for _, d := range req.Details {
		productId, err := util.StringToInt64(d.ProductId)
		if err != nil {
			return nil, err
		}
		details = append(details, &pb.OrderDetailInput{
			ProductId:    productId,
			ProductName:  d.ProductName,
			CategoryType: d.CategoryType,
			Quantity:     d.Quantity,
			UnitPrice:    d.UnitPrice,
			Amount:       d.Amount,
			Remark:       d.Remark,
		})
	}

	supplierId, err := util.StringToInt64(req.SupplierId)
	if err != nil {
		return nil, err
	}
	purchaserId, err := util.StringToInt64(req.PurchaserId)
	if err != nil {
		return nil, err
	}

	no := util.GenerateNo("PO")
	ret, err := l.svcCtx.PurchaseRPC.CreateOrderWithDetails(l.ctx, &pb.CreateOrderWithDetailsReq{
		OrderNo:      no,
		SupplierId:   supplierId,
		OrderDate:    req.OrderDate,
		ExpectedDate: req.ExpectedDate,
		TotalAmount:  req.TotalAmount,
		Status:       req.Status,
		PurchaserId:  purchaserId,
		Details:      details,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.CreateOrderWithDetailsResp{
		OrderId: util.Int64ToString(ret.OrderId),
	}
	return
}
