package purchaseReceipt

import (
	"context"
	hrpb "erp/app/hr/rpc/pb"
	inventorypb "erp/app/inventory/rpc/pb"
	"erp/app/product/rpc/client/product"
	"erp/app/product/rpc/client/productbatch"
	"erp/app/purchase/api/internal/svc"
	"erp/app/purchase/api/internal/types"
	"erp/app/purchase/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetReceiptWithDetailsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetReceiptWithDetailsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetReceiptWithDetailsLogic {
	return &GetReceiptWithDetailsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetReceiptWithDetailsLogic) GetReceiptWithDetails(req *types.GetReceiptWithDetailsReq) (resp *types.GetReceiptWithDetailsResp, err error) {
	receiptId, err := util.StringToInt64(req.ReceiptId)
	if err != nil {
		return nil, err
	}
	ret, err := l.svcCtx.PurchaseRPC.GetReceiptWithDetails(l.ctx, &pb.GetReceiptWithDetailsReq{
		ReceiptId: receiptId,
	})
	if err != nil {
		return nil, err
	}

	// 获取操作人信息
	var createdByNo, createdByName string
	if ret.Receipt.CreatedBy > 0 {
		empResp, err := l.svcCtx.HrRPC.GetEmployeeDetailById(l.ctx, &hrpb.GetEmployeeDetailByIdReq{
			Id: ret.Receipt.CreatedBy,
		})
		if err != nil {
			logx.Errorw("获取操作人信息失败", logx.Field("created_by", ret.Receipt.CreatedBy), logx.Field("error", err))
		} else if empResp.EmployeeNonSensitiveDetail != nil {
			createdByNo = empResp.EmployeeNonSensitiveDetail.EmployeeNo
			createdByName = empResp.EmployeeNonSensitiveDetail.Name
		}
	}

	whResp, err := l.svcCtx.InventoryRPC.GetWarehouseById(l.ctx, &inventorypb.GetWarehouseByIdReq{
		Id: ret.Receipt.WarehouseId,
	})
	if err != nil {
		return nil, err
	}
	warehouseNo := whResp.Warehouse.No
	warehouseName := whResp.Warehouse.Name

	resp = &types.GetReceiptWithDetailsResp{
		Receipt: types.PurchaseReceipt{
			Id:            util.Int64ToString(ret.Receipt.Id),
			ReceiptNo:     ret.Receipt.ReceiptNo,
			OrderId:       util.Int64ToString(ret.Receipt.OrderId),
			WarehouseId:   util.Int64ToString(ret.Receipt.WarehouseId),
			WarehouseNo:   warehouseNo,
			WarehouseName: warehouseName,
			ReceiptDate:   ret.Receipt.ReceiptDate,
			TotalQuantity: ret.Receipt.TotalQuantity,
			TotalAmount:   ret.Receipt.TotalAmount,
			Status:        ret.Receipt.Status,
			CreatedAt:     ret.Receipt.CreatedAt,
			CreatedById:   util.Int64ToString(ret.Receipt.CreatedBy),
			CreatedByNo:   createdByNo,
			CreatedByName: createdByName,
		},
		Details: func() []*types.PurchaseReceiptDetail {
			list := make([]*types.PurchaseReceiptDetail, 0, len(ret.Details))
			productBatchMap := make(map[int64]*productbatch.ProductBatch)
			productMap := make(map[int64]*product.Product)
			for _, d := range ret.Details {
				detail := &types.PurchaseReceiptDetail{
					Id:           util.Int64ToString(d.Id),
					ReceiptId:    util.Int64ToString(d.ReceiptId),
					ProductId:    util.Int64ToString(d.ProductId),
					ProductName:  d.ProductName,
					CategoryType: d.CategoryType,
					Quantity:     d.Quantity,
					UnitPrice:    d.UnitPrice,
					Amount:       d.Amount,
					BatchId:      util.Int64ToString(d.BatchId),
				}
				if d.ProductId > 0 {
					if _, ok := productMap[d.ProductId]; !ok {
						prod, err := l.svcCtx.ProductRPC.GetProductById(l.ctx, &product.GetProductByIdReq{
							Id: d.ProductId,
						})
						if err != nil {
							return nil
						}
						productMap[d.ProductId] = prod.Product
					}
					detail.ProductNo = productMap[d.ProductId].ProductNo
				}
				if d.BatchId > 0 {
					if _, ok := productBatchMap[d.BatchId]; !ok {
						productBatch, err := l.svcCtx.ProductRPC.GetProductBatchById(l.ctx, &productbatch.GetProductBatchByIdReq{
							Id: d.BatchId,
						})
						if err != nil {
							return nil
						}
						productBatchMap[d.BatchId] = productBatch.ProductBatch
					}
					detail.BatchNo = productBatchMap[d.BatchId].BatchNo
				}
				list = append(list, detail)
			}
			return list
		}(),
	}
	return
}
