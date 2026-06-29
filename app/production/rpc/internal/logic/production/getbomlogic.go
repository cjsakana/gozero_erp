package productionlogic

import (
	"context"

	"erp/app/production/rpc/internal/svc"
	"erp/app/production/rpc/production"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBomLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetBomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBomLogic {
	return &GetBomLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetBomLogic) GetBom(in *production.IdReq) (*production.BomInfo, error) {
	// 获取BOM主表
	bom, err := l.svcCtx.BomModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	// 获取BOM明细
	items, err := l.svcCtx.BomItemModel.FindByBomId(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	// 组装返回数据
	var itemInfos []*production.BomItemInfo
	for _, item := range items {
		itemInfos = append(itemInfos, &production.BomItemInfo{
			Id:           item.Id,
			BomId:        item.BomId,
			MaterialId:   item.MaterialId,
			MaterialName: item.MaterialName.String,
			Quantity:     item.Quantity,
			Unit:         item.Unit.String,
			ScrapRate:    item.ScrapRate,
			Remark:       item.Remark.String,
		})
	}

	return &production.BomInfo{
		Id:          bom.Id,
		BomNo:       bom.BomNo,
		ProductId:   bom.ProductId,
		ProductName: bom.ProductName.String,
		Version:     bom.Version,
		UnitCost:    bom.UnitCost,
		IsActive:    bom.IsActive,
		Remark:      bom.Remark.String,
		CreatedAt:   bom.CreatedAt.Unix(),
		CreatedBy:   bom.CreatedBy.Int64,
		UpdatedAt:   bom.UpdatedAt.Unix(),
		UpdatedBy:   bom.UpdatedBy.Int64,
		Items:       itemInfos,
	}, nil
}
