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

type UpdateCustomerCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateCustomerCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCustomerCategoryLogic {
	return &UpdateCustomerCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateCustomerCategoryLogic) UpdateCustomerCategory(in *pb.UpdateCustomerCategoryReq) (*pb.UpdateCustomerCategoryResp, error) {
	err := l.svcCtx.CustomerCategoryModel.XUpdate(l.ctx, &model.CustomerCategory{
		Id:           in.Id,
		Name:         in.Name,
		CreditPolicy: sql.NullString{String: in.CreditPolicy, Valid: true},
	})
	if err != nil {
		return nil, code.UpdateCustomerCategoryFail
	}

	return &pb.UpdateCustomerCategoryResp{}, nil
}
