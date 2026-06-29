package logic

import (
	"context"
	"database/sql"

	"erp/app/customer/rpc/internal/code"
	"erp/app/customer/rpc/internal/svc"
	"erp/app/customer/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
)

type GetCustomerByCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCustomerByCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCustomerByCodeLogic {
	return &GetCustomerByCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCustomerByCodeLogic) GetCustomerByCode(in *pb.GetCustomerByCodeReq) (*pb.GetCustomerByCodeResp, error) {
	one, err := l.svcCtx.CustomerModel.FindOneByCode(l.ctx, sql.NullString{
		String: in.Code,
		Valid:  true,
	})
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, code.CustomerNotFound
		}
		return nil, code.GetCustomerFail
	}

	return &pb.GetCustomerByCodeResp{
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
