package bom

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

type CreateBomLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建BOM
func NewCreateBomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateBomLogic {
	return &CreateBomLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateBomLogic) CreateBom(req *types.CreateBomReq) (resp *types.CreateBomResp, err error) {
	createdBy, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)
	if err != nil {
		return nil, err
	}

	productId, err := util.StringToInt64(req.ProductId)
	if err != nil {
		return nil, err
	}

	productById, err := l.svcCtx.ProductRPC.GetProductById(l.ctx, &product.GetProductByIdReq{
		Id: productId,
	})
	if err != nil {
		return nil, err
	}

	// 调用 RPC 创建 BOM 主表（items 现在通过独立接口创建）
	ret, err := l.svcCtx.ProductionRPC.CreateBom(l.ctx, &production.CreateBomReq{
		ProductId:   productId,
		ProductName: productById.Product.ProductName,
		Version:     req.Version,
		Remark:      req.Remark,
		CreatedBy:   createdBy,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.CreateBomResp{
		Id: util.Int64ToString(ret.Id),
	}
	return
}
