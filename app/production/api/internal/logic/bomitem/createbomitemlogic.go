package bomitem

import (
	"context"
	"erp/app/product/rpc/client/product"

	"erp/app/production/api/internal/svc"
	"erp/app/production/api/internal/types"
	"erp/app/production/rpc/production"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateBomItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建BOM明细
func NewCreateBomItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateBomItemLogic {
	return &CreateBomItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateBomItemLogic) CreateBomItem(req *types.CreateBomItemReq) (resp *types.CreateBomItemResp, err error) {
	// 转换 ID
	bomId, err := util.StringToInt64(req.BomId)
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

	// 调用 RPC 创建 BOM 明细
	_, err = l.svcCtx.ProductionRPC.CreateBomItem(l.ctx, &production.CreateBomItemReq{
		BomId:        bomId,
		MaterialId:   materialId,
		MaterialName: productById.Product.ProductName,
		Quantity:     req.Quantity,
		Unit:         req.Unit,
		ScrapRate:    req.ScrapRate,
		Remark:       req.Remark,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.CreateBomItemResp{}
	return
}
