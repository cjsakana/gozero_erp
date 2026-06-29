package logic

import (
	"context"
	"database/sql"
	"erp/app/customer/rpc/internal/code"
	"erp/app/customer/rpc/internal/model"

	"erp/app/customer/rpc/internal/svc"
	"erp/app/customer/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCustomerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateCustomerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCustomerLogic {
	return &UpdateCustomerLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateCustomerLogic) UpdateCustomer(in *pb.UpdateCustomerReq) (*pb.UpdateCustomerResp, error) {
	err := l.svcCtx.CustomerModel.XUpdate(l.ctx, &model.Customer{
		Id:           in.Id,
		Name:         in.Name,
		CategoryId:   in.CategoryId,
		Contact:      sql.NullString{String: in.Contact, Valid: true},
		Phone:        sql.NullString{String: in.Phone, Valid: true},
		Address:      sql.NullString{String: in.Address, Valid: true},
		CreditLimit:  sql.NullFloat64{Float64: in.CreditLimit, Valid: true},
		UsedCredit:   sql.NullFloat64{Float64: in.UsedCredit, Valid: true},
		PaymentTerms: sql.NullString{String: in.PaymentTerms, Valid: true},
		IsActive:     in.IsActive,
		UpdatedBy:    in.UpdatedBy,
	})
	if err != nil {
		return nil, code.UpdateCustomerFail
	}
	return &pb.UpdateCustomerResp{}, nil
}
