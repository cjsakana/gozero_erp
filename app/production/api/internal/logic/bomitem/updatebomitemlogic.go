package bomitem

import (
	"context"

	"erp/app/production/api/internal/svc"
	"erp/app/production/api/internal/types"
	"erp/app/production/rpc/production"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateBomItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新BOM明细
func NewUpdateBomItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBomItemLogic {
	return &UpdateBomItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateBomItemLogic) UpdateBomItem(req *types.UpdateBomItemReq) (resp *types.UpdateBomItemResp, err error) {
	// 转换 ID
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}

	materialId, err := util.StringToInt64(req.MaterialId)
	if err != nil {
		return nil, err
	}

	// 调用 RPC 更新 BOM 明细
	_, err = l.svcCtx.ProductionRPC.UpdateBomItem(l.ctx, &production.UpdateBomItemReq{
		Id:         id,
		MaterialId: materialId,
		Quantity:   req.Quantity,
		Unit:       req.Unit,
		ScrapRate:  req.ScrapRate,
		Remark:     req.Remark,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.UpdateBomItemResp{}
	return
}
