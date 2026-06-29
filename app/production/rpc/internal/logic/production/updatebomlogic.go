package productionlogic

import (
	"context"
	"database/sql"

	"erp/app/production/rpc/internal/svc"
	"erp/app/production/rpc/production"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateBomLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateBomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBomLogic {
	return &UpdateBomLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateBomLogic) UpdateBom(in *production.UpdateBomReq) (*production.UpdateBomResp, error) {
	// 获取现有BOM
	bom, err := l.svcCtx.BomModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	// 更新BOM主表字段
	if in.Version != "" {
		bom.Version = in.Version
	}
	if in.IsActive > 0 {
		bom.IsActive = int64(in.IsActive)
	}
	bom.Remark = sql.NullString{String: in.Remark, Valid: in.Remark != ""}
	bom.UpdatedBy = sql.NullInt64{Int64: in.UpdatedBy, Valid: in.UpdatedBy > 0}

	// 更新BOM主表
	err = l.svcCtx.BomModel.Update(l.ctx, bom)
	if err != nil {
		return nil, err
	}

	return &production.UpdateBomResp{}, nil
}
