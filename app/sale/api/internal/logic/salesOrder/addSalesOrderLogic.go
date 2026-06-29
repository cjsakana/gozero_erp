package salesOrder

import (
	"context"
	"erp/app/sale/rpc/pb"
	"erp/common/util"
	"erp/common/xtypes"

	"erp/app/sale/api/internal/svc"
	"erp/app/sale/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddSalesOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddSalesOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddSalesOrderLogic {
	return &AddSalesOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddSalesOrderLogic) AddSalesOrder(req *types.AddSalesOrderReq) (resp *types.AddSalesOrderResp, err error) {
	salesmanId, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)
	if err != nil {
		return nil, err
	}

	details := make([]*pb.SalesOrderDetailItem, 0, len(req.Details))
	for _, d := range req.Details {
		productId, err := util.StringToInt64(d.ProductId)
		if err != nil {
			return nil, err
		}
		details = append(details, &pb.SalesOrderDetailItem{
			ProductId:    productId,
			ProductName:  d.ProductName,
			Unit:         d.Unit,
			Quantity:     d.Quantity,
			UnitPrice:    d.UnitPrice,
			Amount:       d.Amount,
			DeliveredQty: 0,
			Remark:       d.Remark,
		})
	}
	customerId, err := util.StringToInt64(req.CustomerId)
	if err != nil {
		return nil, err
	}
	orderNo := util.GenerateNo("SO")
	ret, err := l.svcCtx.SaleRPC.AddSalesOrder(l.ctx, &pb.AddSalesOrderReq{
		OrderNo:      orderNo,
		CustomerId:   customerId,
		OrderDate:    req.OrderDate,
		PromisedDate: req.PromisedDate,
		TotalAmount:  req.TotalAmount,
		Status:       req.Status,
		SalesmanId:   salesmanId,
		ContractUrl:  req.ContractUrl,
		Details:      details,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.AddSalesOrderResp{Id: util.Int64ToString(ret.Id)}
	return
}
