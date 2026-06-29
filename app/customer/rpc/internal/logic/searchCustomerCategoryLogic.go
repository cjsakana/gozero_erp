package logic

import (
	"context"
	"erp/app/customer/rpc/internal/code"
	"erp/app/customer/rpc/internal/svc"
	types2 "erp/app/customer/rpc/internal/types"
	"erp/app/customer/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchCustomerCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchCustomerCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchCustomerCategoryLogic {
	return &SearchCustomerCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchCustomerCategoryLogic) SearchCustomerCategory(in *pb.SearchCustomerCategoryReq) (*pb.SearchCustomerCategoryResp, error) {
	customerCategories, total, err := l.svcCtx.CustomerCategoryModel.Search(l.ctx, &types2.SearchCustomerCategory{
		SearchCom: types2.SearchCom{
			Page:  in.Page,
			Limit: in.Limit,
		},
		Name:         in.Name,
		CreditPolicy: in.CreditPolicy,
	})
	if err != nil {
		return nil, code.GetCustomerFail
	}
	list := make([]*pb.CustomerCategory, total)
	for i, one := range customerCategories {
		list[i] = &pb.CustomerCategory{
			Id:           one.Id,
			Name:         one.Name,
			CreditPolicy: one.CreditPolicy.String,
			CreatedBy:    one.CreatedBy,
			CreatedAt:    one.CreatedAt.Unix(),
		}
	}
	return &pb.SearchCustomerCategoryResp{
		CustomerCategory: list,
		Total:            total,
	}, nil
}
