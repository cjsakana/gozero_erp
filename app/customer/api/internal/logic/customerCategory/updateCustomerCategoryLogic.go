package customerCategory

import (
	"context"
	"erp/app/customer/rpc/customer"
	"erp/common/util"

	"erp/app/customer/api/internal/svc"
	"erp/app/customer/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCustomerCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateCustomerCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCustomerCategoryLogic {
	return &UpdateCustomerCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCustomerCategoryLogic) UpdateCustomerCategory(req *types.UpdateCustomerCategoryReq) (resp *types.UpdateCustomerCategoryResp, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.CustomerRPC.UpdateCustomerCategory(l.ctx, &customer.UpdateCustomerCategoryReq{
		Id:           id,
		Name:         req.Name,
		CreditPolicy: req.CreditPolicy,
	})
	if err != nil {
		return nil, err
	}
	return
}
