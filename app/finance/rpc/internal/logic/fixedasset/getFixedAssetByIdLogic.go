package fixedassetlogic

import (
	"context"
	"erp/app/finance/rpc/internal/code"
	"erp/app/finance/rpc/internal/svc"
	"erp/app/finance/rpc/pb"
	"github.com/zeromicro/go-zero/core/stores/sqlc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFixedAssetByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFixedAssetByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFixedAssetByIdLogic {
	return &GetFixedAssetByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFixedAssetByIdLogic) GetFixedAssetById(in *pb.GetFixedAssetByIdReq) (*pb.GetFixedAssetByIdResp, error) {

	fixedAsset, err := l.svcCtx.FixedAssetModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, code.FixedAssetNotFound
		}
		return nil, code.FixedAssetNotFound
	}

	resp := &pb.GetFixedAssetByIdResp{
		FixedAsset: &pb.FixedAsset{
			Id:                 fixedAsset.Id,
			AssetNo:            fixedAsset.AssetNo.String,
			AssetName:          fixedAsset.AssetName,
			Category:           fixedAsset.Category.String,
			PurchaseDate:       fixedAsset.PurchaseDate.Unix(),
			PurchasePrice:      fixedAsset.PurchasePrice,
			SupplierId:         fixedAsset.SupplierId.Int64,
			DepartmentId:       fixedAsset.DepartmentId,
			UserId:             fixedAsset.UserId.Int64,
			Status:             fixedAsset.Status,
			Location:           fixedAsset.Location.String,
			DepreciationMethod: fixedAsset.DepreciationMethod.String,
			UsefulLife:         fixedAsset.UsefulLife.Int64,
			CreatedAt:          fixedAsset.CreatedAt.Unix(),
		},
	}
	return resp, nil
}
