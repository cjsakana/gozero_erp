package productlogic

import (
	"context"
	"database/sql"
	"erp/app/product/rpc/internal/model"

	"erp/app/product/rpc/internal/svc"
	"erp/app/product/rpc/pb"

	"erp/app/product/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProductLogic {
	return &UpdateProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateProductLogic) UpdateProduct(in *pb.UpdateProductReq) (*pb.UpdateProductResp, error) {
	err := l.svcCtx.ProductModel.XUpdate(l.ctx, &model.Product{
		Id:             in.Id,
		ProductName:    in.ProductName,
		CategoryId:     in.CategoryId,
		Specifications: sql.NullString{String: in.Specifications, Valid: true},
		Unit:           in.Unit,
		PurchasePrice:  sql.NullFloat64{Float64: in.PurchasePrice, Valid: true},
		SellingPrice:   sql.NullFloat64{Float64: in.SellingPrice, Valid: true},
		IsActive:       in.IsActive,
		IsMaterial:     in.IsMaterial,
		UpdatedBy:      in.UpdatedBy,
	})
	if err != nil {

		return nil, code.UpdateProductFail

	}

	return &pb.UpdateProductResp{}, nil
}
