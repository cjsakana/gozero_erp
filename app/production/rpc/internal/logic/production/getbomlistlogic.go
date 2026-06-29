package productionlogic

import (
	"context"
	"erp/app/production/rpc/internal/svc"
	"erp/app/production/rpc/internal/types"
	"erp/app/production/rpc/production"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBomListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetBomListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBomListLogic {
	return &GetBomListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetBomListLogic) GetBomList(in *production.BomListReq) (*production.BomListResp, error) {
	bomList, total, err := l.svcCtx.BomModel.GetBomList(l.ctx, &types.GetBomListParams{
		Page:      in.Page,
		PageSize:  in.PageSize,
		ProductId: in.ProductId,
		IsActive:  in.IsActive,
	})
	if err != nil {
		return nil, err
	}
	var pbBomList []*production.BomInfo

	for _, b := range bomList {
		pbBomList = append(pbBomList, &production.BomInfo{
			Id:          b.Id,
			BomNo:       b.BomNo,
			ProductId:   b.ProductId,
			ProductName: b.ProductName.String,
			Version:     b.Version,
			UnitCost:    b.UnitCost,
			IsActive:    b.IsActive,
			Remark:      b.Remark.String,
			CreatedAt:   b.CreatedAt.Unix(),
			CreatedBy:   b.CreatedBy.Int64,
			UpdatedAt:   b.UpdatedAt.Unix(),
			UpdatedBy:   b.UpdatedBy.Int64,
			Items: func() []*production.BomItemInfo {
				pbItems := []*production.BomItemInfo{}
				items, err := l.svcCtx.BomItemModel.FindByBomId(l.ctx, b.Id)
				if err != nil {
					return nil
				}
				for _, item := range items {
					pbItems = append(pbItems, &production.BomItemInfo{
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
				return pbItems
			}(),
		})
	}

	return &production.BomListResp{
		Total: total,
		List:  pbBomList,
	}, nil
}
