package logic

import (
	"context"
	"erp/app/customer/rpc/internal/code"
	"erp/app/customer/rpc/internal/svc"
	"erp/app/customer/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelCustomerCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelCustomerCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelCustomerCategoryLogic {
	return &DelCustomerCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelCustomerCategoryLogic) DelCustomerCategory(in *pb.DelCustomerCategoryReq) (*pb.DelCustomerCategoryResp, error) {
	err := l.svcCtx.CustomerCategoryModel.XDelete(l.ctx, in.Id)
	if err != nil {
		return nil, code.CustomerCategoryInUse
	}

	return &pb.DelCustomerCategoryResp{}, nil
}
