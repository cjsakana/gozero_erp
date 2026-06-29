package logic

import (
	"context"

	"erp/app/customer/rpc/internal/code"
	"erp/app/customer/rpc/internal/svc"
	"erp/app/customer/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
)

type GetCustomerCategoryByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCustomerCategoryByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCustomerCategoryByIdLogic {
	return &GetCustomerCategoryByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCustomerCategoryByIdLogic) GetCustomerCategoryById(in *pb.GetCustomerCategoryByIdReq) (*pb.GetCustomerCategoryByIdResp, error) {
	one, err := l.svcCtx.CustomerCategoryModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, code.CustomerCategoryNotFound
		}
		return nil, code.GetCustomerCategoryFail
	}

	return &pb.GetCustomerCategoryByIdResp{
		CustomerCategory: &pb.CustomerCategory{
			Id:           one.Id,
			Name:         one.Name,
			CreditPolicy: one.CreditPolicy.String,
			CreatedBy:    one.CreatedBy,
			CreatedAt:    one.CreatedAt.Unix(),
		},
	}, nil
}
