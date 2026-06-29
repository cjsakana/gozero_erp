package productlogic

import (
	"context"
	"erp/app/product/rpc/internal/svc"
	"erp/app/product/rpc/pb"

	"erp/app/product/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
)

type GetProductByNoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProductByNoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductByNoLogic {
	return &GetProductByNoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetProductByNoLogic) GetProductByNo(in *pb.GetProductByNoReq) (*pb.GetProductByNoResp, error) {
	one, err := l.svcCtx.ProductModel.FindOneByProductNo(l.ctx, in.ProductNo)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, code.ProductNotFound
		}
		return nil, code.GetProductFail
	}

	return &pb.GetProductByNoResp{
		Product: &pb.Product{
			Id:             one.Id,
			ProductNo:      one.ProductNo,
			ProductName:    one.ProductName,
			CategoryId:     one.CategoryId,
			Specifications: one.Specifications.String,
			Unit:           one.Unit,
			PurchasePrice:  one.PurchasePrice.Float64,
			SellingPrice:   one.SellingPrice.Float64,
			IsActive:       one.IsActive,
			IsMaterial:     one.IsMaterial,
			CreatedAt:      one.CreatedAt.Unix(),
			CreatedBy:      one.CreatedBy,
			UpdatedAt:      one.UpdatedAt.Unix(),
			UpdatedBy:      one.UpdatedBy,
		},
	}, nil

}
