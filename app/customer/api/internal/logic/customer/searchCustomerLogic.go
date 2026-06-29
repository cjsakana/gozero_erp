package customer

import (
	"context"
	"erp/app/customer/api/internal/svc"
	"erp/app/customer/api/internal/types"
	"erp/app/customer/rpc/customer"
	"erp/app/hr/rpc/client/employeedetail"
	"erp/app/hr/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchCustomerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchCustomerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchCustomerLogic {
	return &SearchCustomerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchCustomerLogic) SearchCustomer(req *types.SearchCustomerReq) (resp *types.SearchCustomerResp, err error) {
	categoryId, err := util.StringToInt64(req.CategoryId)
	if err != nil {
		return nil, err
	}

	ret, err := l.svcCtx.CustomerRPC.SearchCustomer(l.ctx, &customer.SearchCustomerReq{
		Page:         req.Page,
		Limit:        req.Limit,
		Code:         req.Code,
		Uscc:         req.Uscc,
		Name:         req.Name,
		CategoryId:   categoryId,
		Contact:      req.Contact,
		Address:      req.Address,
		PaymentTerms: req.PaymentTerms,
		IsActive:     req.IsActive,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.SearchCustomerResp{
		Total: ret.Total,
	}

	employeeMap := make(map[int64]*employeedetail.EmployeeNonSensitiveDetail)

	for _, v := range ret.Customer {

		if _, ok := employeeMap[v.CreatedBy]; !ok {
			employeeDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb.GetEmployeeDetailByIdReq{
				Id: v.CreatedBy,
			})
			if err != nil {
				logx.Errorf("查询员工信息失败: employeeId=%d, err=%v", v.CreatedBy, err)
				continue
			}
			employeeMap[v.CreatedBy] = employeeDetail.EmployeeNonSensitiveDetail
		}

		if _, ok := employeeMap[v.UpdatedBy]; !ok {
			employeeDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb.GetEmployeeDetailByIdReq{
				Id: v.UpdatedBy,
			})
			if err != nil {
				logx.Errorf("查询员工信息失败: employeeId=%d, err=%v", v.UpdatedBy, err)
				continue
			}
			employeeMap[v.UpdatedBy] = employeeDetail.EmployeeNonSensitiveDetail
		}

		resp.List = append(resp.List, types.Customer{
			Id:            util.Int64ToString(v.Id),
			Code:          v.Code,
			Uscc:          v.Uscc,
			Name:          v.Name,
			CategoryId:    util.Int64ToString(v.CategoryId),
			Contact:       v.Contact,
			Phone:         v.Phone,
			Address:       v.Address,
			CreditLimit:   v.CreditLimit,
			UsedCredit:    v.UsedCredit,
			PaymentTerms:  v.PaymentTerms,
			IsActive:      v.IsActive,
			CreatedBy:     util.Int64ToString(v.CreatedBy),
			CreatedByNo:   employeeMap[v.CreatedBy].EmployeeNo,
			CreatedByName: employeeMap[v.CreatedBy].Name,
			CreatedAt:     v.CreatedAt,
			UpdatedAt:     v.UpdatedAt,
			UpdatedBy:     util.Int64ToString(v.UpdatedBy),
			UpdatedByNo:   employeeMap[v.UpdatedBy].EmployeeNo,
			UpdatedByName: employeeMap[v.UpdatedBy].Name,
		})
	}

	return
}
