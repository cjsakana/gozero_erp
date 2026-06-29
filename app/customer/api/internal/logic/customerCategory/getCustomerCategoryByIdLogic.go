package customerCategory

import (
	"context"
	"erp/app/customer/rpc/customer"
	pb2 "erp/app/hr/rpc/pb"
	"erp/common/util"

	"erp/app/customer/api/internal/svc"
	"erp/app/customer/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCustomerCategoryByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCustomerCategoryByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCustomerCategoryByIdLogic {
	return &GetCustomerCategoryByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCustomerCategoryByIdLogic) GetCustomerCategoryById(req *types.GetCustomerCategoryByIdReq) (resp *types.GetCustomerCategoryByIdResp, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}

	ret, err := l.svcCtx.CustomerRPC.GetCustomerCategoryById(l.ctx, &customer.GetCustomerCategoryByIdReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	empRet, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb2.GetEmployeeDetailByIdReq{
		Id: ret.CustomerCategory.CreatedBy,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.GetCustomerCategoryByIdResp{
		CustomerCategory: types.CustomerCategory{
			Id:            util.Int64ToString(ret.CustomerCategory.Id),
			Name:          ret.CustomerCategory.Name,
			CreditPolicy:  ret.CustomerCategory.CreditPolicy,
			CreatedBy:     util.Int64ToString(ret.CustomerCategory.CreatedBy),
			CreatedByNo:   empRet.EmployeeNonSensitiveDetail.EmployeeNo,
			CreatedByName: empRet.EmployeeNonSensitiveDetail.Name,
			CreatedAt:     ret.CustomerCategory.CreatedAt,
		},
	}

	return
}
