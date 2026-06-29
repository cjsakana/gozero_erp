package customerCategory

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

type SearchCustomerCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchCustomerCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchCustomerCategoryLogic {
	return &SearchCustomerCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchCustomerCategoryLogic) SearchCustomerCategory(req *types.SearchCustomerCategoryReq) (resp *types.SearchCustomerCategoryResp, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}

	ret, err := l.svcCtx.CustomerRPC.SearchCustomerCategory(l.ctx, &customer.SearchCustomerCategoryReq{
		Page:         req.Page,
		Limit:        req.Limit,
		Id:           id,
		Name:         req.Name,
		CreditPolicy: req.CreditPolicy,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.SearchCustomerCategoryResp{
		Total: ret.Total,
	}

	employeeMap := make(map[int64]*employeedetail.EmployeeNonSensitiveDetail)
	for _, v := range ret.CustomerCategory {
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

		resp.List = append(resp.List, types.CustomerCategory{
			Id:            util.Int64ToString(v.Id),
			Name:          v.Name,
			CreditPolicy:  v.CreditPolicy,
			CreatedBy:     util.Int64ToString(v.CreatedBy),
			CreatedByNo:   employeeMap[v.CreatedBy].EmployeeNo,
			CreatedByName: employeeMap[v.CreatedBy].Name,
			CreatedAt:     v.CreatedAt,
		})
	}
	return
}
