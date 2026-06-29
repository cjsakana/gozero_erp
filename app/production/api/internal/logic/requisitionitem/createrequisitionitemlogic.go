package requisitionitem

import (
	"context"
	"erp/app/product/rpc/client/product"

	"erp/app/production/api/internal/svc"
	"erp/app/production/api/internal/types"
	"erp/app/production/rpc/production"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRequisitionItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建领料单明细
func NewCreateRequisitionItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRequisitionItemLogic {
	return &CreateRequisitionItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateRequisitionItemLogic) CreateRequisitionItem(req *types.CreateRequisitionItemReq) (resp *types.CreateRequisitionItemResp, err error) {
	// 转换 ID
	requisitionId, err := util.StringToInt64(req.RequisitionId)
	if err != nil {
		return nil, err
	}

	materialId, err := util.StringToInt64(req.MaterialId)
	if err != nil {
		return nil, err
	}

	productById, err := l.svcCtx.ProductRPC.GetProductById(l.ctx, &product.GetProductByIdReq{
		Id: materialId,
	})
	if err != nil {
		return nil, err
	}

	// 调用 RPC 创建领料单明细
	_, err = l.svcCtx.ProductionRPC.CreateRequisitionItem(l.ctx, &production.CreateRequisitionItemReq{
		RequisitionId:  requisitionId,
		MaterialId:     materialId,
		MaterialName:   productById.Product.ProductName,
		PlanQuantity:   req.PlanQuantity,
		ActualQuantity: req.ActualQuantity,
		Unit:           req.Unit,
		BatchNo:        req.BatchNo,
		Remark:         req.Remark,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.CreateRequisitionItemResp{}
	return
}
