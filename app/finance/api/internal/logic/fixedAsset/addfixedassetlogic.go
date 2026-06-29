package fixedAsset

import (
	"context"
	"erp/app/finance/api/internal/svc"
	"erp/app/finance/api/internal/types"
	"erp/app/finance/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddFixedAssetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddFixedAssetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFixedAssetLogic {
	return &AddFixedAssetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddFixedAssetLogic) AddFixedAsset(req *types.AddFixedAssetReq) (resp *types.AddFixedAssetResp, err error) {
	supplierId, err := util.StringToInt64(req.SupplierId)
	if err != nil {
		return nil, err
	}
	departmentId, err := util.StringToInt64(req.DepartmentId)
	if err != nil {
		return nil, err
	}
	userId, err := util.StringToInt64(req.UserId)
	if err != nil {
		return nil, err
	}
	ret, err := l.svcCtx.FinanceRPC.AddFixedAsset(l.ctx, &pb.AddFixedAssetReq{
		AssetName:          req.AssetName,
		Category:           req.Category,
		PurchaseDate:       req.PurchaseDate,
		PurchasePrice:      req.PurchasePrice,
		SupplierId:         supplierId,
		DepartmentId:       departmentId,
		UserId:             userId,
		Status:             req.Status,
		Location:           req.Location,
		DepreciationMethod: req.DepreciationMethod,
		UsefulLife:         req.UsefulLife,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.AddFixedAssetResp{
		Id: util.Int64ToString(ret.Id),
	}
	return
}
