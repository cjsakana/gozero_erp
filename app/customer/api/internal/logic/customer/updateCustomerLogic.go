package customer

import (
	"context"
	"erp/app/customer/rpc/customer"
	"erp/common/util"
	"erp/common/xtypes"

	"erp/app/customer/api/internal/svc"
	"erp/app/customer/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCustomerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateCustomerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCustomerLogic {
	return &UpdateCustomerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCustomerLogic) UpdateCustomer(req *types.UpdateCustomerReq) (resp *types.UpdateCustomerResp, err error) {
	employeeId, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)
	if err != nil {
		return nil, err
	}
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	categoryId, err := util.StringToInt64(req.CategoryId)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.CustomerRPC.UpdateCustomer(l.ctx, &customer.UpdateCustomerReq{
		Id:           id,
		Name:         req.Name,
		CategoryId:   categoryId,
		Contact:      req.Contact,
		Phone:        req.Phone,
		Address:      req.Address,
		CreditLimit:  req.CreditLimit,
		UsedCredit:   req.UsedCredit,
		PaymentTerms: req.PaymentTerms,
		IsActive:     req.IsActive,
		UpdatedBy:    employeeId,
	})
	if err != nil {
		return nil, err
	}

	return
}
