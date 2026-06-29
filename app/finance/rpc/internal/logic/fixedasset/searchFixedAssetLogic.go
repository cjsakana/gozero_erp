package fixedassetlogic

import (
	"context"

	"erp/app/finance/rpc/internal/svc"
	"erp/app/finance/rpc/pb"

	"erp/app/finance/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type SearchFixedAssetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchFixedAssetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchFixedAssetLogic {
	return &SearchFixedAssetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchFixedAssetLogic) SearchFixedAsset(in *pb.SearchFixedAssetReq) (*pb.SearchFixedAssetResp, error) {
	fixedAssets, total, err := l.svcCtx.FixedAssetModel.Search(
		l.ctx,
		in.AssetNo,
		in.AssetName,
		in.Category,
		in.SupplierId,
		in.DepartmentId,
		in.Status,
		in.Page,
		in.Limit,
	)
	if err != nil {

		return nil, code.GetFixedAssetFail

	}

	var pbFixedAssets []*pb.FixedAsset
	for _, fa := range fixedAssets {
		pbFixedAssets = append(pbFixedAssets, &pb.FixedAsset{
			Id:                 fa.Id,
			AssetNo:            fa.AssetNo.String,
			AssetName:          fa.AssetName,
			Category:           fa.Category.String,
			PurchaseDate:       fa.PurchaseDate.Unix(),
			PurchasePrice:      fa.PurchasePrice,
			SupplierId:         fa.SupplierId.Int64,
			DepartmentId:       fa.DepartmentId,
			UserId:             fa.UserId.Int64,
			Status:             fa.Status,
			Location:           fa.Location.String,
			DepreciationMethod: fa.DepreciationMethod.String,
			UsefulLife:         fa.UsefulLife.Int64,
			CreatedAt:          fa.CreatedAt.Unix(),
		})
	}

	return &pb.SearchFixedAssetResp{
		FixedAsset: pbFixedAssets,
		Total:      total,
	}, nil
}
