package customer

import (
	"context"
	"erp/app/customer/api/internal/svc"
	"erp/app/customer/api/internal/types"
	"erp/app/customer/rpc/pb"
	pb2 "erp/app/hr/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCustomerByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCustomerByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCustomerByIdLogic {
	return &GetCustomerByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCustomerByIdLogic) GetCustomerById(req *types.GetCustomerByIdReq) (resp *types.GetCustomerByIdResp, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	
	ret, err := l.svcCtx.CustomerRPC.GetCustomerById(l.ctx, &pb.GetCustomerByIdReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	empRet, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb2.GetEmployeeDetailByIdReq{
		Id: ret.Customer.CreatedBy,
	})
	if err != nil {
		return nil, err
	}

	updatedByNo := empRet.EmployeeNonSensitiveDetail.EmployeeNo
	updatedByName := empRet.EmployeeNonSensitiveDetail.Name

	// 创建和修改的不是同一人
	if ret.Customer.CreatedBy != ret.Customer.UpdatedBy {
		empRet2, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb2.GetEmployeeDetailByIdReq{
			Id: ret.Customer.UpdatedBy,
		})
		if err != nil {
			return nil, err
		}
		updatedByNo = empRet2.EmployeeNonSensitiveDetail.EmployeeNo
		updatedByName = empRet2.EmployeeNonSensitiveDetail.Name
	}

	resp = &types.GetCustomerByIdResp{
		Customer: types.Customer{
			Id:            util.Int64ToString(ret.Customer.Id),
			Code:          ret.Customer.Code,
			Uscc:          ret.Customer.Uscc,
			Name:          ret.Customer.Name,
			CategoryId:    util.Int64ToString(ret.Customer.CategoryId),
			Contact:       ret.Customer.Contact,
			Phone:         ret.Customer.Phone,
			Address:       ret.Customer.Address,
			CreditLimit:   ret.Customer.CreditLimit,
			UsedCredit:    ret.Customer.UsedCredit,
			PaymentTerms:  ret.Customer.PaymentTerms,
			IsActive:      ret.Customer.IsActive,
			CreatedBy:     util.Int64ToString(ret.Customer.CreatedBy),
			CreatedByNo:   empRet.EmployeeNonSensitiveDetail.EmployeeNo,
			CreatedByName: empRet.EmployeeNonSensitiveDetail.Name,
			CreatedAt:     ret.Customer.CreatedAt,
			UpdatedAt:     ret.Customer.UpdatedAt,
			UpdatedBy:     util.Int64ToString(ret.Customer.UpdatedBy),
			UpdatedByNo:   updatedByNo,
			UpdatedByName: updatedByName,
		},
	}

	return
}
