package fixedassetlogic

import (
	"context"
	"database/sql"
	"erp/common/util"
	"time"

	"erp/app/finance/rpc/internal/model"
	"erp/app/finance/rpc/internal/svc"
	"erp/app/finance/rpc/pb"

	"erp/app/finance/rpc/internal/code"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddFixedAssetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddFixedAssetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFixedAssetLogic {
	return &AddFixedAssetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------fixedAsset-----------------------
func (l *AddFixedAssetLogic) AddFixedAsset(in *pb.AddFixedAssetReq) (*pb.AddFixedAssetResp, error) {
	id := util.GenerateSnowflake()
	assetNo := util.GenerateNo("FA")
	data := &model.FixedAsset{
		Id:                 id,
		AssetNo:            sql.NullString{String: assetNo, Valid: true},
		AssetName:          in.AssetName,
		Category:           sql.NullString{String: in.Category, Valid: in.Category != ""},
		PurchaseDate:       time.Unix(in.PurchaseDate, 0),
		PurchasePrice:      in.PurchasePrice,
		SupplierId:         sql.NullInt64{Int64: in.SupplierId, Valid: in.SupplierId != 0},
		DepartmentId:       in.DepartmentId,
		UserId:             sql.NullInt64{Int64: in.UserId, Valid: in.UserId != 0},
		Status:             in.Status,
		Location:           sql.NullString{String: in.Location, Valid: in.Location != ""},
		DepreciationMethod: sql.NullString{String: in.DepreciationMethod, Valid: in.DepreciationMethod != ""},
		UsefulLife:         sql.NullInt64{Int64: in.UsefulLife, Valid: in.UsefulLife != 0},
	}

	_, err := l.svcCtx.FixedAssetModel.Insert(l.ctx, data)
	if err != nil {

		return nil, code.AddFixedAssetFail

	}

	return &pb.AddFixedAssetResp{Id: id}, nil
}
