package purchaseorderlogic

import (
	"context"

	"erp/app/purchase/rpc/internal/svc"
	"erp/app/purchase/rpc/internal/types"
	"erp/app/purchase/rpc/pb"

	"erp/app/purchase/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type SearchOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchOrderLogic {
	return &SearchOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询采购订单
func (l *SearchOrderLogic) SearchOrder(in *pb.SearchOrderReq) (*pb.SearchOrderResp, error) {
	params := &types.SearchOrderParams{
		SearchComm: types.SearchComm{Page: in.Page, Limit: in.Limit},
		OrderNo:    in.OrderNo,
		SupplierId: in.SupplierId,
		Status:     in.Status,
	}
	orders, total, err := l.svcCtx.PurchaseOrderModel.Search(l.ctx, params)
	if err != nil {

		return nil, code.GetOrderFail

	}

	var pbOrdersD []*pb.PurchaseOrderWithDetails
	for _, o := range orders {
		pbOrderWithDetails := &pb.PurchaseOrderWithDetails{
			Orders: &pb.PurchaseOrder{
				Id:         o.Id,
				OrderNo:    o.OrderNo,
				SupplierId: o.SupplierId,
				OrderDate:  o.OrderDate.Unix(),
				ExpectedDate: func() int64 {
					if o.ExpectedDate.Valid {
						return o.ExpectedDate.Time.Unix()
					}
					return 0
				}(),
				TotalAmount: o.TotalAmount,
				Status:      o.Status,
				PurchaserId: o.PurchaserId,
				ContractUrl: o.ContractUrl.String,
				CreatedAt:   o.CreatedAt.Unix(),
				UpdatedAt:   o.UpdatedAt.Unix(),
			},
		}
		od, err := l.svcCtx.PurchaseOrderDetailModel.ListByOrderId(l.ctx, o.Id)
		if err != nil {

			return nil, code.GetOrderFail

		}
		pbOrderWithDetails.Total = int64(len(od))
		for _, detail := range od {
			pbOrderWithDetails.Details = append(pbOrderWithDetails.Details, &pb.PurchaseOrderDetail{
				Id:           detail.Id,
				OrderId:      detail.OrderId,
				ProductId:    detail.ProductId.Int64,
				ProductName:  detail.ProductName.String,
				CategoryType: detail.CategoryType,
				Quantity:     detail.Quantity,
				UnitPrice:    detail.UnitPrice,
				Amount:       detail.Amount,
				ReceivedQty:  detail.ReceivedQty,
				Remark:       detail.Remark.String,
			})
		}

		pbOrdersD = append(pbOrdersD, pbOrderWithDetails)
	}

	return &pb.SearchOrderResp{
		OrdersWithDetails: pbOrdersD,
		Total:             total,
	}, nil
}
