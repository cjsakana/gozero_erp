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

type CreateRequisitionItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateRequisitionItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRequisitionItemLogic {
	return &CreateRequisitionItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 领料单明细管理
func (l *CreateRequisitionItemLogic) CreateRequisitionItem(in *production.CreateRequisitionItemReq) (*production.CreateRequisitionItemResp, error) {
	// 验证领料单是否存在
	_, err := l.svcCtx.MaterialRequisitionModel.FindOne(l.ctx, in.RequisitionId)
	if err != nil {
		return nil, err
	}

	id := util.GenerateSnowflake()

	// 插入领料单明细
	_, err = l.svcCtx.MaterialRequisitionItemModel.Insert(l.ctx, &model.MaterialRequisitionItem{
		Id:             id,
		RequisitionId:  in.RequisitionId,
		MaterialId:     in.MaterialId,
		MaterialName:   sql.NullString{String: in.MaterialName, Valid: true},
		PlanQuantity:   in.PlanQuantity,
		ActualQuantity: in.ActualQuantity,
		Unit:           sql.NullString{String: in.Unit, Valid: in.Unit != ""},
		BatchNo:        sql.NullString{String: in.BatchNo, Valid: in.BatchNo != ""},
		Remark:         sql.NullString{String: in.Remark, Valid: in.Remark != ""},
	})
	if err != nil {
		return nil, err
	}

	return &production.CreateRequisitionItemResp{Id: id}, nil
}
