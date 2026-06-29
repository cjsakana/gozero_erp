package purchaseReceipt

import (
	"context"
	"erp/app/hr/rpc/client/employeedetail"
	"erp/app/inventory/rpc/client/inventory"
	"erp/app/purchase/api/internal/svc"
	"erp/app/purchase/api/internal/types"
	"erp/app/purchase/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchReceiptLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchReceiptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchReceiptLogic {
	return &SearchReceiptLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchReceiptLogic) SearchReceipt(req *types.SearchReceiptReq) (resp *types.SearchReceiptResp, err error) {
	var orderId, warehouseId int64
	if req.OrderId != "" {
		orderId, err = util.StringToInt64(req.OrderId)
		if err != nil {
			return nil, err
		}
	}
	if req.WarehouseId != "" {
		warehouseId, err = util.StringToInt64(req.WarehouseId)
		if err != nil {
			return nil, err
		}
	}

	ret, err := l.svcCtx.PurchaseRPC.SearchReceipt(l.ctx, &pb.SearchReceiptReq{
		Page:        req.Page,
		Limit:       req.Limit,
		ReceiptNo:   req.ReceiptNo,
		OrderId:     orderId,
		WarehouseId: warehouseId,
	})
	if err != nil {
		return nil, err
	}

	warehouseMap := make(map[int64]*inventory.WarehouseDetail)
	//productBatchMap := make(map[int64]*productbatch.ProductBatch)
	//productMap := make(map[int64]*product.Product)
	employeeMap := make(map[int64]*employeedetail.EmployeeNonSensitiveDetail)

	var receiptWithDetails []*types.PurchaseReceiptWithDetails
	for _, rd := range ret.ReceiptsWithDetails {
		// 获取仓库和操作人信息
		if _, ok := warehouseMap[rd.Receipts.WarehouseId]; !ok {
			w, e := l.svcCtx.InventoryRPC.GetWarehouseById(l.ctx, &inventory.GetWarehouseByIdReq{
				Id: rd.Receipts.WarehouseId,
			})
			if e == nil {
				warehouseMap[rd.Receipts.WarehouseId] = w.Warehouse
			}
		}

		if _, ok := employeeMap[rd.Receipts.CreatedBy]; !ok {
			empResp, err := l.svcCtx.HrRPC.GetEmployeeDetailById(l.ctx, &employeedetail.GetEmployeeDetailByIdReq{
				Id: rd.Receipts.CreatedBy,
			})
			if err == nil {
				employeeMap[rd.Receipts.CreatedBy] = empResp.EmployeeNonSensitiveDetail
			}
		}

		receiptWithDetails = append(receiptWithDetails, &types.PurchaseReceiptWithDetails{
			Receipt: types.PurchaseReceipt{
				Id:            util.Int64ToString(rd.Receipts.Id),
				ReceiptNo:     rd.Receipts.ReceiptNo,
				OrderId:       util.Int64ToString(rd.Receipts.OrderId),
				WarehouseId:   util.Int64ToString(rd.Receipts.WarehouseId),
				WarehouseNo:   warehouseMap[rd.Receipts.WarehouseId].No,
				WarehouseName: warehouseMap[rd.Receipts.WarehouseId].Name,
				ReceiptDate:   rd.Receipts.ReceiptDate,
				TotalQuantity: rd.Receipts.TotalQuantity,
				TotalAmount:   rd.Receipts.TotalAmount,
				Status:        rd.Receipts.Status,
				CreatedAt:     rd.Receipts.CreatedAt,
				CreatedById:   util.Int64ToString(rd.Receipts.CreatedBy),
				CreatedByNo:   employeeMap[rd.Receipts.CreatedBy].EmployeeNo,
				CreatedByName: employeeMap[rd.Receipts.CreatedBy].Name,
			},
			Total: rd.Total,
			//Details: func() []*types.PurchaseReceiptDetail {
			//	details := make([]*types.PurchaseReceiptDetail, 0, rd.Total)
			//	for _, d := range rd.Details {
			//
			//		// 商品、商品批次
			//		if _, ok := productMap[d.ProductId]; !ok {
			//			prod, err := l.svcCtx.ProductRPC.GetProductById(l.ctx, &product.GetProductByIdReq{
			//				Id: d.ProductId,
			//			})
			//			if err == nil {
			//				productMap[d.ProductId] = prod.Product
			//			}
			//		}
			//
			//		if _, ok := productBatchMap[d.BatchId]; !ok {
			//			productBatch, err := l.svcCtx.ProductRPC.GetProductBatchById(l.ctx, &productbatch.GetProductBatchByIdReq{
			//				Id: d.BatchId,
			//			})
			//			if err == nil {
			//				productBatchMap[d.BatchId] = productBatch.ProductBatch
			//			}
			//		}
			//
			//		details = append(details, &types.PurchaseReceiptDetail{
			//			Id:           util.Int64ToString(d.Id),
			//			ReceiptId:    util.Int64ToString(d.ReceiptId),
			//			ProductId:    util.Int64ToString(d.ProductId),
			//			ProductNo:    productMap[d.ProductId].ProductNo,
			//			ProductName:  d.ProductName,
			//			CategoryType: d.CategoryType,
			//			Quantity:     d.Quantity,
			//			UnitPrice:    d.UnitPrice,
			//			Amount:       d.Amount,
			//			BatchId:      util.Int64ToString(d.BatchId),
			//			BatchNo:      productBatchMap[d.BatchId].BatchNo,
			//		})
			//	}
			//	return details
			//}(),
			Details: nil,
		})
	}

	resp = &types.SearchReceiptResp{
		ReceiptsWithDetails: receiptWithDetails,
		Total:               ret.Total,
	}
	return
}
