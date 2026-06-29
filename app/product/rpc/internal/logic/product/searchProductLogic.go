package productlogic

import (
	"context"
	"erp/app/product/rpc/internal/code"
	"erp/app/product/rpc/internal/svc"
	types2 "erp/app/product/rpc/internal/types"
	"erp/app/product/rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type SearchProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchProductLogic {
	return &SearchProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchProductLogic) SearchProduct(in *pb.SearchProductReq) (*pb.SearchProductResp, error) {
	products, total, err := l.svcCtx.ProductModel.Search(l.ctx, &types2.SearchProductParams{
		SearchCom: types2.SearchCom{
			Page:  in.Page,
			Limit: in.Limit,
		},
		ProductNo:   in.ProductNo,
		ProductName: in.ProductName,
		CategoryId:  in.CategoryId,
		IsActive:    in.IsActive,
		IsMaterial:  in.IsMaterial,
	})
	if err != nil {
		return nil, code.GetProductFail

	}

	productList := make([]*pb.Product, 0)
	for _, one := range products {
		productList = append(productList, &pb.Product{
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
		})
	}
	return &pb.SearchProductResp{
		Total:   total,
		Product: productList,
	}, nil
}
