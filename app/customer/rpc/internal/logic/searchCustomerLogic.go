package logic

import (
	"context"
	"erp/app/customer/rpc/internal/code"
	"erp/app/customer/rpc/internal/svc"
	"erp/app/customer/rpc/internal/types"
	"erp/app/customer/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchCustomerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchCustomerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchCustomerLogic {
	return &SearchCustomerLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchCustomerLogic) SearchCustomer(in *pb.SearchCustomerReq) (*pb.SearchCustomerResp, error) {
	customers, total, err := l.svcCtx.CustomerModel.Search(l.ctx, &types.SearchCustomer{
		SearchCom: types.SearchCom{
			Page:  in.Page,
			Limit: in.Limit,
		},
		Code:         in.Code,
		USCC:         in.Uscc,
		Name:         in.Name,
		CategoryId:   in.CategoryId,
		Contact:      in.Contact,
		Address:      in.Address,
		PaymentTerms: in.PaymentTerms,
		IsActive:     in.IsActive,
	})
	if err != nil {
		return nil, code.GetCustomerCategoryFail
	}

	list := make([]*pb.Customer, total)
	for i, one := range customers {
		list[i] = &pb.Customer{
			Id:           one.Id,
			Code:         one.Code.String,
			Uscc:         one.Uscc.String,
			Name:         one.Name,
			CategoryId:   one.CategoryId,
			Contact:      one.Contact.String,
			Phone:        one.Phone.String,
			Address:      one.Address.String,
			CreditLimit:  one.CreditLimit.Float64,
			UsedCredit:   one.UsedCredit.Float64,
			PaymentTerms: one.PaymentTerms.String,
			IsActive:     one.IsActive,
			CreatedAt:    one.CreatedAt.Unix(),
			CreatedBy:    one.CreatedBy,
			UpdatedAt:    one.UpdatedAt.Unix(),
			UpdatedBy:    one.UpdatedBy,
		}
	}

	return &pb.SearchCustomerResp{
		Customer: list,
		Total:    total,
	}, nil
}
