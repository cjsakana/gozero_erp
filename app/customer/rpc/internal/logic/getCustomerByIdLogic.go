package logic

import (
	"context"

	"erp/app/customer/rpc/internal/code"
	"erp/app/customer/rpc/internal/svc"
	"erp/app/customer/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
)

type GetCustomerByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCustomerByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCustomerByIdLogic {
	return &GetCustomerByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCustomerByIdLogic) GetCustomerById(in *pb.GetCustomerByIdReq) (*pb.GetCustomerByIdResp, error) {
	one, err := l.svcCtx.CustomerModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, code.CustomerNotFound
		}
		return nil, code.GetCustomerFail
	}

	return &pb.GetCustomerByIdResp{
		Customer: &pb.Customer{
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
		},
	}, nil
}
