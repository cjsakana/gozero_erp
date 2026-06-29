package satisfaction

import (
	"context"
	"erp/app/customer/rpc/pb"
	"erp/common/util"
	"erp/common/xtypes"

	"erp/app/customer/api/internal/svc"
	"erp/app/customer/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddCustomerSatisfactionSurveyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddCustomerSatisfactionSurveyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCustomerSatisfactionSurveyLogic {
	return &AddCustomerSatisfactionSurveyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddCustomerSatisfactionSurveyLogic) AddCustomerSatisfactionSurvey(req *types.AddCustomerSatisfactionSurveyReq) (resp *types.AddCustomerSatisfactionSurveyResp, err error) {
	createdBy, err := util.GetInt64FromCtx(l.ctx, xtypes.EmployeeIdKey)
	if err != nil {
		return nil, err
	}

	ret, err := l.svcCtx.CustomerRPC.GetCustomerByCode(l.ctx, &pb.GetCustomerByCodeReq{Code: req.Code})
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.CustomerRPC.AddCustomerSatisfactionSurvey(l.ctx, &pb.AddCustomerSatisfactionSurveyReq{
		CustomerId:    ret.Customer.Id,
		QualityScore:  req.QualityScore,
		DeliveryScore: req.DeliveryScore,
		ServiceScore:  req.ServiceScore,
		OverallScore:  req.OverallScore,
		Remark:        req.Remark,
		CreatedBy:     createdBy,
	})
	if err != nil {
		return nil, err
	}

	return
}
