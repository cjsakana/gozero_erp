package logic

import (
	"context"
	"erp/app/supplier/rpc/internal/code"
	"erp/app/supplier/rpc/internal/svc"
	types2 "erp/app/supplier/rpc/internal/types"
	"erp/app/supplier/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchSupplierLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchSupplierLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchSupplierLogic {
	return &SearchSupplierLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchSupplierLogic) SearchSupplier(in *pb.SearchSupplierReq) (*pb.SearchSupplierResp, error) {
	suppliers, total, err := l.svcCtx.SupplierModel.Search(l.ctx, &types2.SearchSupplierParams{
		SearchCom: types2.SearchCom{
			Page:  in.Page,
			Limit: in.Limit,
		},
		Code:         in.Code,
		Uscc:         in.Uscc,
		Name:         in.Name,
		Contact:      in.Contact,
		Address:      in.Address,
		PaymentTerms: in.PaymentTerms,
		Credit:       in.Credit,
		IsActive:     in.IsActive,
	})
	if err != nil {
		return nil, code.SupplierEvaluationNotFound
	}
	list := []*pb.Supplier{}
	for _, supplier := range suppliers {
		list = append(list, &pb.Supplier{
			Id:           supplier.Id,
			Code:         supplier.Code.String,
			Uscc:         supplier.Uscc.String,
			Name:         supplier.Name,
			Contact:      supplier.Contact.String,
			Phone:        supplier.Phone.String,
			Address:      supplier.Address.String,
			PaymentTerms: supplier.PaymentTerms.String,
			Credit:       supplier.Credit.String,
			IsActive:     supplier.IsActive,
			CreatedAt:    supplier.CreatedAt.Unix(),
			CreatedBy:    supplier.CreatedBy,
			UpdatedAt:    supplier.UpdatedAt.Unix(),
			UpdatedBy:    supplier.UpdatedBy,
		})
	}

	return &pb.SearchSupplierResp{
		Supplier: list,
		Total:    total,
	}, nil
}
