package bom

import (
	"context"
	"erp/app/production/rpc/production"
	"erp/common/util"
	"erp/common/xtypes"

	"erp/app/production/api/internal/svc"
	"erp/app/production/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateBomLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新BOM
func NewUpdateBomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBomLogic {
	return &UpdateBomLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateBomLogic) UpdateBom(req *types.UpdateBomReq) (resp *types.UpdateBomResp, err error) {
	updatedBy, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)
	if err != nil {
		return nil, err
	}
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}

	// 调用 RPC 更新 BOM 主表（items 现在通过独立接口管理）
	_, err = l.svcCtx.ProductionRPC.UpdateBom(l.ctx, &production.UpdateBomReq{
		Id:        id,
		Version:   req.Version,
		IsActive:  req.IsActive,
		Remark:    req.Remark,
		UpdatedBy: updatedBy,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.UpdateBomResp{}
	return
}
