package salesOrder

import (
	"context"
	"erp/app/customer/rpc/customer"
	"erp/app/hr/rpc/client/employeedetail"
	"erp/app/sale/api/internal/svc"
	"erp/app/sale/api/internal/types"
	"erp/app/sale/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchSalesOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchSalesOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchSalesOrderLogic {
	return &SearchSalesOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchSalesOrderLogic) SearchSalesOrder(req *types.SearchSalesOrderReq) (resp *types.SearchSalesOrderResp, err error) {
	customerId, err := util.StringToInt64(req.CustomerId)
	if err != nil {
		return nil, err
	}
	ret, err := l.svcCtx.SaleRPC.SearchSalesOrder(l.ctx, &pb.SearchSalesOrderReq{
		Page:           req.Page,
		Limit:          req.Limit,
		OrderNo:        req.OrderNo,
		CustomerId:     customerId,
		Status:         req.Status,
		StartOrderDate: req.StartOrderDate,
		EndOrderDate:   req.EndOrderDate,
	})
	if err != nil {
		return nil, err
	}

	salesmanMap := make(map[int64]*employeedetail.EmployeeNonSensitiveDetail)
	customerMap := make(map[int64]*customer.Customer)

	// 组装响应列表
	orderList := make([]types.OrderWithDetails, 0, len(ret.OrderWithDetails))
	for _, o := range ret.OrderWithDetails {
		if _, ok := salesmanMap[o.SalesOrder.SalesmanId]; !ok {
			employeeDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &employeedetail.GetEmployeeDetailByIdReq{
				Id: o.SalesOrder.SalesmanId,
			})
			if err != nil {
				logx.Errorf("查询员工信息失败: employeeId=%d, err=%v", o.SalesOrder.SalesmanId, err)
			}
			salesmanMap[o.SalesOrder.SalesmanId] = employeeDetail.EmployeeNonSensitiveDetail
		}

		if _, ok := customerMap[o.SalesOrder.CustomerId]; !ok {
			customerInfo, err := l.svcCtx.CustomerRPC.GetCustomerById(l.ctx, &customer.GetCustomerByIdReq{
				Id: o.SalesOrder.CustomerId,
			})
			if err != nil {
				logx.Errorf("查询客户信息失败: CustomerId=%d, err=%v", o.SalesOrder.CustomerId, err)
			}
			customerMap[o.SalesOrder.CustomerId] = customerInfo.Customer
		}

		salesOrder := types.SalesOrder{
			Id:           util.Int64ToString(o.SalesOrder.Id),
			OrderNo:      o.SalesOrder.OrderNo,
			CustomerId:   util.Int64ToString(o.SalesOrder.CustomerId),
			CustomerCode: customerMap[o.SalesOrder.CustomerId].Code,
			CustomerName: customerMap[o.SalesOrder.CustomerId].Name,
			OrderDate:    o.SalesOrder.OrderDate,
			PromisedDate: o.SalesOrder.PromisedDate,
			TotalAmount:  o.SalesOrder.TotalAmount,
			Status:       o.SalesOrder.Status,
			SalesmanId:   util.Int64ToString(o.SalesOrder.SalesmanId),
			SalesmanNo:   salesmanMap[o.SalesOrder.SalesmanId].EmployeeNo,
			SalesmanName: salesmanMap[o.SalesOrder.SalesmanId].Name,
			ContractUrl:  o.SalesOrder.ContractUrl,
			CreatedAt:    o.SalesOrder.CreatedAt,
		}

		orderList = append(orderList, types.OrderWithDetails{
			SalesOrder:       salesOrder,
			Total:            o.Total,
			SalesOrderDetail: nil,
			//SalesOrderDetail: func() []types.SalesOrderDetail {
			//	detailList := make([]types.SalesOrderDetail, 0)
			//	for _, d := range o.SalesOrderDetail {
			//		detailList = append(detailList, types.SalesOrderDetail{
			//			Id:           util.Int64ToString(d.Id),
			//			OrderId:      util.Int64ToString(d.OrderId),
			//			ProductId:    util.Int64ToString(d.ProductId),
			//			ProductNo:    "",
			//			ProductName:  d.ProductName,
			//			Unit:         d.Unit,
			//			Quantity:     d.Quantity,
			//			UnitPrice:    d.UnitPrice,
			//			Amount:       d.Amount,
			//			DeliveredQty: 0,
			//			Remark:       d.Remark,
			//		})
			//	}
			//	return detailList
			//}(),

		})
	}

	resp = &types.SearchSalesOrderResp{
		Total:            ret.Total,
		OrderWithDetails: orderList,
	}
	return
}
