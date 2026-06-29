package customerCategory

import (
	"context"
	"erp/app/customer/rpc/customer"
	"erp/common/util"

	"erp/app/customer/api/internal/svc"
	"erp/app/customer/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelCustomerCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelCustomerCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelCustomerCategoryLogic {
	return &DelCustomerCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelCustomerCategoryLogic) DelCustomerCategory(req *types.DelCustomerCategoryReq) (resp *types.DelCustomerCategoryResp, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.CustomerRPC.DelCustomerCategory(l.ctx, &customer.DelCustomerCategoryReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return
}
