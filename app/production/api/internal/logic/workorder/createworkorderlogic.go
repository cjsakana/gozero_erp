package workorder

import (
	"context"
	"erp/app/product/rpc/client/product"
	"erp/app/production/rpc/production"
	"erp/common/util"
	"erp/common/xtypes"

	"erp/app/production/api/internal/svc"
	"erp/app/production/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateWorkOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建生产工单
func NewCreateWorkOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateWorkOrderLogic {
	return &CreateWorkOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateWorkOrderLogic) CreateWorkOrder(req *types.CreateWorkOrderReq) (resp *types.CreateWorkOrderResp, err error) {
	createdBy, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)
	if err != nil {
		return nil, err
	}
	productId, err := util.StringToInt64(req.ProductId)
	if err != nil {
		return nil, err
	}
	bomId, err := util.StringToInt64(req.BomId)
	if err != nil {
		return nil, err
	}
	warehouseId, err := util.StringToInt64(req.WarehouseId)
	if err != nil {
		return nil, err
	}
	productById, err := l.svcCtx.ProductRPC.GetProductById(l.ctx, &product.GetProductByIdReq{
		Id: productId,
	})
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.ProductionRPC.CreateWorkOrder(l.ctx, &production.CreateWorkOrderReq{
		ProductId:     productId,
		ProductName:   productById.Product.ProductName,
		BomId:         bomId,
		Quantity:      req.Quantity,
		WarehouseId:   warehouseId,
		Priority:      req.Priority,
		PlanStartDate: req.PlanStartDate,
		PlanEndDate:   req.PlanEndDate,
		Workshop:      req.Workshop,
		Remark:        req.Remark,
		CreatedBy:     createdBy,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.CreateWorkOrderResp{}
	return
}
