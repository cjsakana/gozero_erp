package productlogic

import (
	"context"
	"database/sql"
	"erp/app/product/rpc/internal/model"
	"erp/app/product/rpc/internal/svc"
	"erp/app/product/rpc/pb"
	"erp/common/util"

	"erp/app/product/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddProductLogic {
	return &AddProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------product-----------------------
func (l *AddProductLogic) AddProduct(in *pb.AddProductReq) (*pb.AddProductResp, error) {
	id := util.GenerateSnowflake()
	no := ""
	if in.CategoryId == 1000000000000000001 {
		//原材料
		no = util.GenerateNo("MAT")
	} else if in.CategoryId == 1000000000000000002 {
		// 标准件
		no = util.GenerateNo("STD")
	} else if in.CategoryId == 1000000000000000003 {
		// 成品
		no = util.GenerateNo("PRD")
	} else if in.CategoryId == 1000000000000000004 {
		// 半成品
		no = util.GenerateNo("SEM")
	} else {
		// 生产辅助
		no = util.GenerateNo("AUX")
	}

	_, err := l.svcCtx.ProductModel.Insert(l.ctx, &model.Product{
		Id:             id,
		ProductNo:      no,
		ProductName:    in.ProductName,
		CategoryId:     in.CategoryId,
		Specifications: sql.NullString{String: in.Specifications, Valid: true},
		Unit:           in.Unit,
		PurchasePrice:  sql.NullFloat64{Float64: in.PurchasePrice, Valid: true},
		SellingPrice:   sql.NullFloat64{Float64: in.SellingPrice, Valid: true},
		IsActive:       in.IsActive,
		IsMaterial:     in.IsMaterial,
		CreatedBy:      in.CreatedBy,
		UpdatedBy:      in.CreatedBy,
	})
	if err != nil {
		return nil, code.AddProductFail
	}

	return &pb.AddProductResp{
		Id: id,
	}, nil
}
