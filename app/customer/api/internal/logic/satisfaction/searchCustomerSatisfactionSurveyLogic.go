package satisfaction

import (
	"context"
	"erp/app/customer/api/internal/svc"
	"erp/app/customer/api/internal/types"
	"erp/app/customer/rpc/customer"
	"erp/app/hr/rpc/client/employeedetail"
	"erp/common/util"
	"github.com/zeromicro/go-zero/core/logx"
)

type SearchCustomerSatisfactionSurveyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchCustomerSatisfactionSurveyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchCustomerSatisfactionSurveyLogic {
	return &SearchCustomerSatisfactionSurveyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchCustomerSatisfactionSurveyLogic) SearchCustomerSatisfactionSurvey(req *types.SearchCustomerSatisfactionSurveyReq) (resp *types.SearchCustomerSatisfactionSurveyResp, err error) {
	customerId, err := util.StringToInt64(req.CustomerId)
	if err != nil {
		return nil, err
	}

	ret, err := l.svcCtx.CustomerRPC.SearchCustomerSatisfactionSurvey(l.ctx, &customer.SearchCustomerSatisfactionSurveyReq{
		Page:        req.Page,
		Limit:       req.Limit,
		CustomerId:  customerId,
		QualityMin:  req.QualityMin,
		QualityMax:  req.QualityMax,
		QualityOp:   req.QualityOp,
		DeliveryMin: req.DeliveryMin,
		DeliveryMax: req.DeliveryMax,
		DeliveryOp:  req.DeliveryOp,
		ServiceMin:  req.ServiceMin,
		ServiceMax:  req.ServiceMax,
		ServiceOp:   req.ServiceOp,
		OverallMin:  req.OverallMin,
		OverallMax:  req.OverallMax,
		OverallOp:   req.OverallOp,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.SearchCustomerSatisfactionSurveyResp{
		Total: ret.Total,
	}

	employeeMap := make(map[int64]*employeedetail.EmployeeNonSensitiveDetail)
	for _, v := range ret.CustomerSatisfactionSurvey {
		if _, ok := employeeMap[v.CreatedBy]; !ok {
			empRet, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &employeedetail.GetEmployeeDetailByIdReq{
				Id: v.CreatedBy,
			})
			if err != nil {
				continue
			}
			employeeMap[v.CreatedBy] = empRet.EmployeeNonSensitiveDetail
		}

		resp.List = append(resp.List, types.CustomerSatisfactionSurvey{
			Id:            util.Int64ToString(v.Id),
			CustomerId:    util.Int64ToString(v.CustomerId),
			QualityScore:  v.QualityScore,
			DeliveryScore: v.DeliveryScore,
			ServiceScore:  v.ServiceScore,
			OverallScore:  v.OverallScore,
			Remark:        v.Remark,
			CreatedAt:     v.CreatedAt,
			CreatedBy:     util.Int64ToString(v.CreatedBy),
			CreatedByNo:   employeeMap[v.CreatedBy].EmployeeNo,
			CreatedByName: employeeMap[v.CreatedBy].Name,
		})
	}

	return
}
