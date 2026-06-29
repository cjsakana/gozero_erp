package inventory

import (
	"context"
	"erp/app/hr/rpc/client/employeedetail"
	"erp/app/inventory/api/internal/svc"
	"erp/app/inventory/api/internal/types"
	"erp/app/inventory/rpc/client/inventory"
	"erp/app/inventory/rpc/pb"
	"erp/app/product/rpc/client/productbatch"
	pb2 "erp/app/product/rpc/pb"
	"erp/common/util"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetInventoryTransactionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetInventoryTransactionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetInventoryTransactionsLogic {
	return &GetInventoryTransactionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetInventoryTransactionsLogic) GetInventoryTransactions(req *types.GetInventoryTransactionsReq) (resp *types.GetInventoryTransactionsResp, err error) {
	// 初始化返回对象
	resp = &types.GetInventoryTransactionsResp{}

	productId, err := util.StringToInt64(req.ProductId)
	if err != nil {
		return nil, err
	}
	warehouseId, err := util.StringToInt64(req.WarehouseId)
	if err != nil {
		return nil, err
	}

	// 初始化缓存 map
	warehouseMap := make(map[int64]*inventory.WarehouseDetail)
	productMap := make(map[int64]*pb2.Product)
	batchMap := make(map[int64]*productbatch.ProductBatch)
	employeeMap := make(map[int64]*employeedetail.EmployeeNonSensitiveDetail)

	// 如果提供了产品名称，先查询对应的产品ID列表

	// 查询库存交易记录
	var allTransactions []*pb.InventoryTransaction
	var total int64

	// 直接查询，不再在这里过滤 productId（RPC 层不支持批量 ID 查询）
	ret, err := l.svcCtx.InventoryRPC.SearchInventoryTransaction(l.ctx, &pb.SearchInventoryTransactionReq{
		Page:            req.Page,
		Limit:           req.Limit,
		ProductId:       productId,
		WarehouseId:     warehouseId,
		TransactionType: req.TransactionType,
		ReferenceType:   req.ReferenceType,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
	})
	if err != nil {
		return nil, err
	}
	allTransactions = ret.InventoryTransaction
	total = ret.Total
	resp = &types.GetInventoryTransactionsResp{
		Total: total,
	}
	for _, transaction := range allTransactions {
		// 获取产品信息
		if _, ok := productMap[transaction.ProductId]; !ok {
			ret4, err := l.svcCtx.ProductRPC.GetProductById(l.ctx, &pb2.GetProductByIdReq{
				Id: transaction.ProductId,
			})
			if err != nil {
				return nil, err
			}
			productMap[transaction.ProductId] = ret4.Product
		}

		// 获取仓库信息
		if _, ok := warehouseMap[transaction.WarehouseId]; !ok {
			ret3, err := l.svcCtx.InventoryRPC.GetWarehouseById(l.ctx, &inventory.GetWarehouseByIdReq{
				Id: transaction.WarehouseId,
			})
			if err != nil {
				return nil, err
			}
			warehouseMap[transaction.WarehouseId] = ret3.Warehouse
		}

		// 获取批次信息（如果有）
		var batchNo string
		if transaction.BatchId > 0 {
			if _, ok := batchMap[transaction.BatchId]; !ok {
				batchResp, err := l.svcCtx.ProductRPC.GetProductBatchById(l.ctx, &pb2.GetProductBatchByIdReq{
					Id: transaction.BatchId,
				})
				if err != nil {
					return nil, err
				}
				batchMap[transaction.BatchId] = batchResp.ProductBatch
			}
			batchNo = batchMap[transaction.BatchId].BatchNo
		}

		// 获取操作员信息（如果有）
		var operatorNo, operatorName string
		if transaction.OperatorId > 0 {
			if _, ok := employeeMap[transaction.OperatorId]; !ok {
				empResp, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &employeedetail.GetEmployeeDetailByIdReq{
					Id: transaction.OperatorId,
				})
				if err != nil {
					return nil, err
				}
				employeeMap[transaction.OperatorId] = empResp.EmployeeNonSensitiveDetail
			}
			operatorNo = employeeMap[transaction.OperatorId].EmployeeNo
			operatorName = employeeMap[transaction.OperatorId].Name
		}

		resp.Transactions = append(resp.Transactions, types.InventoryTransaction{
			Id:              util.Int64ToString(transaction.Id),
			ProductId:       util.Int64ToString(transaction.ProductId),
			ProductNo:       productMap[transaction.ProductId].ProductNo,
			ProductName:     productMap[transaction.ProductId].ProductName,
			WarehouseId:     util.Int64ToString(transaction.WarehouseId),
			WarehouseNo:     warehouseMap[transaction.WarehouseId].No,
			WarehouseName:   warehouseMap[transaction.WarehouseId].Name,
			BatchId:         util.Int64ToString(transaction.BatchId),
			BatchNo:         batchNo,
			TransactionType: transaction.TransactionType,
			Quantity:        transaction.Quantity,
			ReferenceType:   transaction.ReferenceType,
			ReferenceId:     transaction.ReferenceId,
			OperatorId:      util.Int64ToString(transaction.OperatorId),
			OperatorNo:      operatorNo,
			OperatorName:    operatorName,
			CreatedAt:       transaction.CreatedAt,
		})
	}
	resp.Total = int64(len(resp.Transactions))
	return
}
