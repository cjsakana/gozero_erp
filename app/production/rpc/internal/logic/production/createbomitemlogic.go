package productionlogic

import (
	"context"
	"database/sql"
	"erp/common/util"

	"erp/app/production/rpc/internal/model"
	"erp/app/production/rpc/internal/svc"
	"erp/app/production/rpc/production"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateBomItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateBomItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateBomItemLogic {
	return &CreateBomItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BOM明细管理
func (l *CreateBomItemLogic) CreateBomItem(in *production.CreateBomItemReq) (*production.CreateBomItemResp, error) {
	// 验证 BOM 是否存在
	_, err := l.svcCtx.BomModel.FindOne(l.ctx, in.BomId)
	if err != nil {
		return nil, err
	}

	id := util.GenerateSnowflake()

	// 插入 BOM 明细
	_, err = l.svcCtx.BomItemModel.Insert(l.ctx, &model.BomItem{
		Id:           id,
		BomId:        in.BomId,
		MaterialId:   in.MaterialId,
		MaterialName: sql.NullString{String: in.MaterialName, Valid: true},
		Quantity:     in.Quantity,
		Unit:         sql.NullString{String: in.Unit, Valid: in.Unit != ""},
		ScrapRate:    in.ScrapRate,
		Remark:       sql.NullString{String: in.Remark, Valid: in.Remark != ""},
	})
	if err != nil {
		return nil, err
	}

	return &production.CreateBomItemResp{Id: id}, nil
}
