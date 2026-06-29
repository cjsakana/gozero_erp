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

type AddCustomerCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddCustomerCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCustomerCategoryLogic {
	return &AddCustomerCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------customerCategory-----------------------
func (l *AddCustomerCategoryLogic) AddCustomerCategory(in *pb.AddCustomerCategoryReq) (*pb.AddCustomerCategoryResp, error) {
	id := util.GenerateSnowflake()
	_, err := l.svcCtx.CustomerCategoryModel.Insert(l.ctx, &model.CustomerCategory{
		Id:           id,
		Name:         in.Name,
		CreditPolicy: sql.NullString{String: in.CreditPolicy, Valid: true},
		CreatedBy:    in.CreatedBy,
	})
	if err != nil {
		return nil, code.AddCustomerCategoryFail
	}
	return &pb.AddCustomerCategoryResp{
		Id: id,
	}, nil
}
