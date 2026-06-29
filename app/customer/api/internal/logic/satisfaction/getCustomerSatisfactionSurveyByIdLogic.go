package satisfaction

import (
	"context"
	"erp/app/customer/rpc/pb"
	pb2 "erp/app/hr/rpc/pb"
	"erp/common/util"

	"erp/app/customer/api/internal/svc"
	"erp/app/customer/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCustomerSatisfactionSurveyByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCustomerSatisfactionSurveyByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCustomerSatisfactionSurveyByIdLogic {
	return &GetCustomerSatisfactionSurveyByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCustomerSatisfactionSurveyByIdLogic) GetCustomerSatisfactionSurveyById(req *types.GetCustomerSatisfactionSurveyByIdReq) (resp *types.GetCustomerSatisfactionSurveyByIdResp, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}

	ret, err := l.svcCtx.CustomerRPC.GetCustomerSatisfactionSurveyById(l.ctx, &pb.GetCustomerSatisfactionSurveyByIdReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	empRet, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb2.GetEmployeeDetailByIdReq{
		Id: ret.CustomerSatisfactionSurvey.CreatedBy,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.GetCustomerSatisfactionSurveyByIdResp{
		CustomerSatisfactionSurvey: types.CustomerSatisfactionSurvey{
			Id:            util.Int64ToString(ret.CustomerSatisfactionSurvey.Id),
			CustomerId:    util.Int64ToString(ret.CustomerSatisfactionSurvey.CustomerId),
			QualityScore:  ret.CustomerSatisfactionSurvey.QualityScore,
			DeliveryScore: ret.CustomerSatisfactionSurvey.DeliveryScore,
			ServiceScore:  ret.CustomerSatisfactionSurvey.ServiceScore,
			OverallScore:  ret.CustomerSatisfactionSurvey.OverallScore,
			Remark:        ret.CustomerSatisfactionSurvey.Remark,
			CreatedAt:     ret.CustomerSatisfactionSurvey.CreatedAt,
			CreatedBy:     util.Int64ToString(ret.CustomerSatisfactionSurvey.CreatedBy),
			CreatedByNo:   empRet.EmployeeNonSensitiveDetail.EmployeeNo,
			CreatedByName: empRet.EmployeeNonSensitiveDetail.Name,
		},
	}

	return
}
