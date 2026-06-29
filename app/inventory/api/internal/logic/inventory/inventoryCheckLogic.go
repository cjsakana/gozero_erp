package inventory

import (
	"context"
	"erp/app/inventory/api/internal/svc"
	"erp/app/inventory/api/internal/types"
	"erp/app/inventory/rpc/pb"
	pb2 "erp/app/product/rpc/pb"
	"erp/common/util"
	"erp/common/xtypes"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
)

type InventoryCheckLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInventoryCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InventoryCheckLogic {
	return &InventoryCheckLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InventoryCheckLogic) InventoryCheck(req *types.InventoryCheckReq) (resp *types.InventoryCheckResp, err error) {
	var successCount int64
	var failCount int64
	var checkDetails []*types.InventoryCheckDetail

	warehouseId, err := util.StringToInt64(req.WarehouseId)
	if err != nil {
		return nil, err
	}

	// 验证仓库是否存在
	_, err = l.svcCtx.InventoryRPC.GetWarehouseById(l.ctx, &pb.GetWarehouseByIdReq{
		Id: warehouseId,
	})
	if err != nil {
		return nil, fmt.Errorf("仓库不存在: %v", err)
	}

	// 从上下文获取当前用户ID
	employeeId, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)
	if err != nil {
		return nil, err
	}

	// 处理每个盘点项
	for _, item := range req.CheckItems {
		detail, err := l.processCheckItem(l.ctx, warehouseId, item, employeeId)
		if err != nil {
			detail.Success = false
			detail.Message = err.Error()
			failCount++
		} else {
			detail.Success = true
			detail.Message = "盘点成功"
			successCount++
		}
		checkDetails = append(checkDetails, detail)
	}

	return &types.InventoryCheckResp{
		SuccessCount: successCount,
		FailCount:    failCount,
		CheckDetails: checkDetails,
	}, nil
}

// 处理单个盘点项
func (l *InventoryCheckLogic) processCheckItem(ctx context.Context, warehouseId int64, item types.InventoryCheckItem, operatorId int64) (*types.InventoryCheckDetail, error) {
	detail := &types.InventoryCheckDetail{
		ProductId:   item.ProductId,
		ActualStock: item.ActualStock,
	}

	productId, err := util.StringToInt64(item.ProductId)
	if err != nil {
		return nil, err
	}
	batchId, err := util.StringToInt64(item.BatchId)
	if err != nil {
		return nil, err
	}

	// 1. 获取商品信息
	ret, err := l.svcCtx.ProductRPC.GetProductById(ctx, &pb2.GetProductByIdReq{
		Id: productId,
	})
	if err != nil {
		detail.Message = "商品不存在"
		return detail, fmt.Errorf("商品不存在")
	}
	detail.ProductName = ret.Product.ProductName

	// 2. 获取当前库存
	ret2, err := l.svcCtx.InventoryRPC.SearchInventory(ctx, &pb.SearchInventoryReq{
		ProductId:   productId,
		WarehouseId: warehouseId,
		Limit:       1,
	})
	if err != nil {
		detail.Message = "获取库存失败"
		return detail, fmt.Errorf("获取库存失败: %v", err)
	}

	// 如果库存记录不存在，创建新记录
	if len(ret2.Inventory) == 0 {
		detail.BeforeStock = 0
		detail.Difference = item.ActualStock
		_, err := l.svcCtx.InventoryRPC.AddInventory(l.ctx, &pb.AddInventoryReq{
			ProductId:       productId,
			WarehouseId:     warehouseId,
			CurrentStock:    item.ActualStock,
			SafetyStock:     0,
			LockedStock:     0,
			BatchId:         batchId,
			TransactionType: 3, // 盘点调整
			ReferenceType:   4, // 盘盈
			ReferenceId:     "Inventory Check",
			OperatorId:      operatorId,
		})
		if err != nil {
			return nil, err
		}
		return detail, nil
	}

	inventoryE := ret2.Inventory[0]
	detail.BeforeStock = inventoryE.CurrentStock
	detail.Difference = item.ActualStock - inventoryE.CurrentStock

	// 3. 更新库存
	transactionType := 0
	if item.ActualStock > inventoryE.CurrentStock {
		transactionType = 4 // 盘盈
	} else {
		transactionType = 5 // 盘亏
	}
	_, err = l.svcCtx.InventoryRPC.UpdateInventory(l.ctx, &pb.UpdateInventoryReq{
		InventoryId:     inventoryE.InventoryId,
		ProductId:       productId,
		WarehouseId:     warehouseId,
		CurrentStock:    item.ActualStock,
		AdjustType:      3, // set
		TransactionType: 3, // 盘点调整
		ReferenceType:   int64(transactionType),
		ReferenceId:     "Inventory Check",
		BatchId:         batchId,
		OperatorId:      operatorId,
	})
	if err != nil {
		return nil, err
	}
	return detail, nil
}
