package productionlogic

import (
	"context"
	"database/sql"
	"erp/app/production/rpc/internal/model"
	"erp/app/production/rpc/internal/svc"
	"erp/app/production/rpc/production"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateBomLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateBomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateBomLogic {
	return &CreateBomLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BOM管理
func (l *CreateBomLogic) CreateBom(in *production.CreateBomReq) (*production.CreateBomResp, error) {
	// 生成BOM主表雪花ID
	bomId := util.GenerateSnowflake()

	bomNo := util.GenerateNo("BOM")

	// 创建BOM主表（items 现在通过独立接口创建）
	bom := &model.Bom{
		Id:          bomId,
		BomNo:       bomNo,
		ProductId:   in.ProductId,
		ProductName: sql.NullString{String: in.ProductName, Valid: true},
		Version:     in.Version,
		UnitCost:    0,
		IsActive:    1,
		Remark:      sql.NullString{String: in.Remark, Valid: in.Remark != ""},
		CreatedBy:   sql.NullInt64{Int64: in.CreatedBy, Valid: in.CreatedBy > 0},
		UpdatedBy:   sql.NullInt64{},
	}

	_, err := l.svcCtx.BomModel.Insert(l.ctx, bom)
	if err != nil {
		return nil, err
	}
	return &production.CreateBomResp{Id: bomId}, nil
}
