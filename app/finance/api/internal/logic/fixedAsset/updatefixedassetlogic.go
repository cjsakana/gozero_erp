package fixedAsset

import (
	"context"
	"erp/app/finance/api/internal/svc"
	"erp/app/finance/api/internal/types"
	"erp/app/finance/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFixedAssetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateFixedAssetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFixedAssetLogic {
	return &UpdateFixedAssetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateFixedAssetLogic) UpdateFixedAsset(req *types.UpdateFixedAssetReq) (resp *types.UpdateFixedAssetResp, err error) {
	id, err := util.StringToInt64(req.Id)
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
	_, err = l.svcCtx.FinanceRPC.UpdateFixedAsset(l.ctx, &pb.UpdateFixedAssetReq{
		Id:                 id,
		AssetName:          req.AssetName,
		Category:           req.Category,
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

	resp = &types.UpdateFixedAssetResp{}
	return
}
