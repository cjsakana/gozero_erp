package customer

import (
	"context"
	"erp/app/customer/api/internal/svc"
	"erp/app/customer/api/internal/types"
	"erp/app/customer/rpc/customer"
	"erp/common/util"
	"erp/common/xtypes"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddCustomerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddCustomerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCustomerLogic {
	return &AddCustomerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddCustomerLogic) AddCustomer(req *types.AddCustomerReq) (resp *types.AddCustomerResp, err error) {
	createdBy, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)
	if err != nil {
		return nil, err
	}

	categoryId, err := util.StringToInt64(req.CategoryId)
	if err != nil {
		return nil, err
	}

	code := util.GenerateNo("CUST")

	ret, err := l.svcCtx.CustomerRPC.AddCustomer(l.ctx, &customer.AddCustomerReq{
		Code:         code,
		Uscc:         req.Uscc,
		Name:         req.Name,
		CategoryId:   categoryId,
		Contact:      req.Contact,
		Phone:        req.Phone,
		Address:      req.Address,
		CreditLimit:  req.CreditLimit,
		UsedCredit:   req.UsedCredit,
		PaymentTerms: req.PaymentTerms,
		IsActive:     req.IsActive,
		CreatedBy:    createdBy,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.AddCustomerResp{
		Id: util.Int64ToString(ret.Id),
	}

	return
}
