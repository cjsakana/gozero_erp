package logic

import (
	"context"
	"database/sql"
	"erp/app/customer/rpc/internal/code"
	"erp/app/customer/rpc/internal/model"
	"erp/common/util"

	"erp/app/customer/rpc/internal/svc"
	"erp/app/customer/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddCustomerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddCustomerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCustomerLogic {
	return &AddCustomerLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------customer-----------------------
func (l *AddCustomerLogic) AddCustomer(in *pb.AddCustomerReq) (*pb.AddCustomerResp, error) {
	id := util.GenerateSnowflake()
	_, err := l.svcCtx.CustomerModel.Insert(l.ctx, &model.Customer{
		Id:           id,
		Code:         sql.NullString{String: in.Code, Valid: true},
		Uscc:         sql.NullString{String: in.Uscc, Valid: true},
		Name:         in.Name,
		CategoryId:   in.CategoryId,
		Contact:      sql.NullString{String: in.Contact, Valid: true},
		Phone:        sql.NullString{String: in.Phone, Valid: true},
		Address:      sql.NullString{String: in.Address, Valid: true},
		CreditLimit:  sql.NullFloat64{Float64: in.CreditLimit, Valid: true},
		UsedCredit:   sql.NullFloat64{Float64: in.UsedCredit, Valid: true},
		PaymentTerms: sql.NullString{String: in.PaymentTerms, Valid: true},
		IsActive:     in.IsActive,
		CreatedBy:    in.CreatedBy,
		UpdatedBy:    in.CreatedBy, //记录最后一次更新的，创建也是更新
	})
	if err != nil {
		return nil, code.AddCustomerFail
	}

	return &pb.AddCustomerResp{
		Id: id,
	}, nil
}
