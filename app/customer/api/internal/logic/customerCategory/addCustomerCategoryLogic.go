package customerCategory

import (
	"context"
	"erp/app/customer/rpc/customer"
	"erp/common/util"
	"erp/common/xtypes"

	"erp/app/customer/api/internal/svc"
	"erp/app/customer/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddCustomerCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddCustomerCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCustomerCategoryLogic {
	return &AddCustomerCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddCustomerCategoryLogic) AddCustomerCategory(req *types.AddCustomerCategoryReq) (resp *types.AddCustomerCategoryResp, err error) {
	createdBy, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)
	if err != nil {
		return nil, err
	}
	ret, err := l.svcCtx.CustomerRPC.AddCustomerCategory(l.ctx, &customer.AddCustomerCategoryReq{
		Name:         req.Name,
		CreditPolicy: req.CreditPolicy,
		CreatedBy:    createdBy,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.AddCustomerCategoryResp{
		Id: util.Int64ToString(ret.Id),
	}
	return
}
