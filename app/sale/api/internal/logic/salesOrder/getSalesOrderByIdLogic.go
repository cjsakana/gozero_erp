package salesOrder

import (
	"context"
	"erp/app/customer/rpc/customer"
	"erp/app/hr/rpc/client/employeedetail"
	"erp/app/product/rpc/client/product"
	"erp/app/sale/api/internal/svc"
	"erp/app/sale/api/internal/types"
	"erp/app/sale/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSalesOrderByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSalesOrderByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSalesOrderByIdLogic {
	return &GetSalesOrderByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSalesOrderByIdLogic) GetSalesOrderById(req *types.GetSalesOrderByIdReq) (resp *types.GetSalesOrderByIdResp, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	ret, err := l.svcCtx.SaleRPC.GetSalesOrderById(l.ctx, &pb.GetSalesOrderByIdReq{Id: id})
	if err != nil {
		return nil, err
	}

	// 获取销售员信息
	var salesmanNo, salesmanName string
	if ret.OrderWithDetails.SalesOrder.SalesmanId > 0 {
		empResp, err := l.svcCtx.HrRPC.GetEmployeeDetailById(l.ctx, &employeedetail.GetEmployeeDetailByIdReq{
			Id: ret.OrderWithDetails.SalesOrder.SalesmanId,
		})
		if err != nil {
			logx.Errorw("获取销售员信息失败", logx.Field("salesman_id", ret.OrderWithDetails.SalesOrder.SalesmanId), logx.Field("error", err))
		} else if empResp.EmployeeNonSensitiveDetail != nil {
			salesmanNo = empResp.EmployeeNonSensitiveDetail.EmployeeNo
			salesmanName = empResp.EmployeeNonSensitiveDetail.Name
		}
	}
	var customerCode, customerName string
	if ret.OrderWithDetails.SalesOrder.CustomerId > 0 {
		customerInfo, err := l.svcCtx.CustomerRPC.GetCustomerById(l.ctx, &customer.GetCustomerByIdReq{
			Id: ret.OrderWithDetails.SalesOrder.CustomerId,
		})
		if err != nil {
			logx.Errorw("获取客户信息失败", logx.Field("CustomerId", ret.OrderWithDetails.SalesOrder.SalesmanId), logx.Field("error", err))
		} else if customerInfo.Customer != nil {
			customerCode = customerInfo.Customer.Code
			customerName = customerInfo.Customer.Name
		}
	}

	resp = &types.GetSalesOrderByIdResp{
		OrderWithDetails: types.OrderWithDetails{
			SalesOrder: types.SalesOrder{
				Id:           util.Int64ToString(ret.OrderWithDetails.SalesOrder.Id),
				OrderNo:      ret.OrderWithDetails.SalesOrder.OrderNo,
				CustomerId:   util.Int64ToString(ret.OrderWithDetails.SalesOrder.CustomerId),
				CustomerCode: customerCode,
				CustomerName: customerName,
				OrderDate:    ret.OrderWithDetails.SalesOrder.OrderDate,
				PromisedDate: ret.OrderWithDetails.SalesOrder.PromisedDate,
				TotalAmount:  ret.OrderWithDetails.SalesOrder.TotalAmount,
				Status:       ret.OrderWithDetails.SalesOrder.Status,
				SalesmanId:   util.Int64ToString(ret.OrderWithDetails.SalesOrder.SalesmanId),
				SalesmanNo:   salesmanNo,
				SalesmanName: salesmanName,
				ContractUrl:  ret.OrderWithDetails.SalesOrder.ContractUrl,
				CreatedAt:    ret.OrderWithDetails.SalesOrder.CreatedAt,
			},
			Total: ret.OrderWithDetails.Total,
			SalesOrderDetail: func() []types.SalesOrderDetail {
				list := make([]types.SalesOrderDetail, 0, len(ret.OrderWithDetails.SalesOrderDetail))

				productMap := make(map[int64]*product.Product)

				for _, d := range ret.OrderWithDetails.SalesOrderDetail {

					if _, ok := productMap[d.ProductId]; !ok {
						prod, err := l.svcCtx.ProductRPC.GetProductById(l.ctx, &product.GetProductByIdReq{
							Id: d.ProductId,
						})
						if err != nil {
							return nil
						}
						productMap[d.ProductId] = prod.Product
					}

					list = append(list, types.SalesOrderDetail{
						Id:           util.Int64ToString(d.Id),
						OrderId:      util.Int64ToString(d.OrderId),
						ProductId:    util.Int64ToString(d.ProductId),
						ProductNo:    productMap[d.ProductId].ProductNo,
						ProductName:  d.ProductName,
						Unit:         d.Unit,
						Quantity:     d.Quantity,
						UnitPrice:    d.UnitPrice,
						Amount:       d.Amount,
						DeliveredQty: d.DeliveredQty,
						Remark:       d.Remark,
					})
				}
				return list
			}(),
		},
	}
	return
}
