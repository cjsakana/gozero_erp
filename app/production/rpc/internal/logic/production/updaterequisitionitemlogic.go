package productionlogic

import (
	"context"
	"database/sql"
	"fmt"

	"erp/app/production/rpc/internal/svc"
	"erp/app/production/rpc/internal/types"
	"erp/app/production/rpc/production"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRequisitionItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateRequisitionItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRequisitionItemLogic {
	return &UpdateRequisitionItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateRequisitionItemLogic) UpdateRequisitionItem(in *production.UpdateRequisitionItemReq) (*production.UpdateRequisitionItemResp, error) {
	// 查找现有记录
	requisitionItem, err := l.svcCtx.MaterialRequisitionItemModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	requisitionItem.MaterialId = in.MaterialId
	requisitionItem.PlanQuantity = in.PlanQuantity
	requisitionItem.ActualQuantity = in.ActualQuantity
	requisitionItem.Unit = sql.NullString{String: in.Unit, Valid: in.Unit != ""}
	requisitionItem.BatchNo = sql.NullString{String: in.BatchNo, Valid: in.BatchNo != ""}
	requisitionItem.Remark = sql.NullString{String: in.Remark, Valid: in.Remark != ""}

	// 更新数据库
	err = l.svcCtx.MaterialRequisitionItemModel.Update(l.ctx, requisitionItem)
	if err != nil {
		return nil, err
	}

	// 清除缓存
	cacheKey := fmt.Sprintf("%s%v", types.CacheErpProductionMaterialRequisitionItemIdPrefix, in.Id)
	l.svcCtx.BizRedis.DelCtx(l.ctx, cacheKey)

	return &production.UpdateRequisitionItemResp{}, nil
}
