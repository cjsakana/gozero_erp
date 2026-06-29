package fixedassetlogic

import (
	"context"
	"database/sql"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"time"

	"erp/app/finance/rpc/internal/model"
	"erp/app/finance/rpc/internal/svc"
	"erp/app/finance/rpc/pb"

	"erp/app/finance/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFixedAssetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateFixedAssetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFixedAssetLogic {
	return &UpdateFixedAssetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateFixedAssetLogic) UpdateFixedAsset(in *pb.UpdateFixedAssetReq) (*pb.UpdateFixedAssetResp, error) {

	fixedAsset, err := l.svcCtx.FixedAssetModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, code.FixedAssetNotFound
		}
		return nil, code.FixedAssetNotFound
	}
	data := &model.FixedAsset{
		Id:                 fixedAsset.Id,
		AssetNo:            fixedAsset.AssetNo,
		AssetName:          in.AssetName,
		Category:           sql.NullString{String: in.Category, Valid: true},
		PurchaseDate:       time.Unix(in.PurchaseDate, 0),
		PurchasePrice:      in.PurchasePrice,
		SupplierId:         sql.NullInt64{Int64: in.SupplierId, Valid: in.SupplierId != 0},
		DepartmentId:       in.DepartmentId,
		UserId:             sql.NullInt64{Int64: in.UserId, Valid: in.UserId != 0},
		Status:             in.Status,
		Location:           sql.NullString{String: in.Location, Valid: true},
		DepreciationMethod: sql.NullString{String: in.DepreciationMethod, Valid: true},
		UsefulLife:         sql.NullInt64{Int64: in.UsefulLife, Valid: in.UsefulLife != 0},
		CreatedAt:          time.Time{},
	}

	err = l.svcCtx.FixedAssetModel.Update(l.ctx, data)
	if err != nil {

		return nil, code.UpdateFixedAssetFail

	}

	return &pb.UpdateFixedAssetResp{}, nil
}
