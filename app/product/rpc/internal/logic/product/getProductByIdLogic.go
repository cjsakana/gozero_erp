package productlogic

import (
	"context"
	"erp/app/product/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/stores/sqlc"

	"erp/app/product/rpc/internal/svc"
	"erp/app/product/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProductByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductByIdLogic {
	return &GetProductByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetProductByIdLogic) GetProductById(in *pb.GetProductByIdReq) (*pb.GetProductByIdResp, error) {
	one, err := l.svcCtx.ProductModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, code.ProductNotFound
		}
		return nil, code.GetProductFail
	}

	return &pb.GetProductByIdResp{
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
