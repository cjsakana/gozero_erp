package logic

import (
	"context"
	"database/sql"
	"erp/app/customer/rpc/internal/code"
	"erp/app/customer/rpc/internal/model"
	"erp/app/customer/rpc/internal/svc"
	"erp/app/customer/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddCustomerSatisfactionSurveyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddCustomerSatisfactionSurveyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCustomerSatisfactionSurveyLogic {
	return &AddCustomerSatisfactionSurveyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------customerSatisfactionSurvey-----------------------
func (l *AddCustomerSatisfactionSurveyLogic) AddCustomerSatisfactionSurvey(in *pb.AddCustomerSatisfactionSurveyReq) (*pb.AddCustomerSatisfactionSurveyResp, error) {
	id := util.GenerateSnowflake()
	_, err := l.svcCtx.SatisfactionModel.Insert(l.ctx, &model.CustomerSatisfactionSurvey{
		Id:            id,
		CustomerId:    in.CustomerId,
		QualityScore:  in.QualityScore,
		DeliveryScore: in.DeliveryScore,
		ServiceScore:  in.ServiceScore,
		OverallScore:  in.OverallScore,
		Remark:        sql.NullString{String: in.Remark, Valid: true},
		CreatedBy:     in.CreatedBy,
	})
	if err != nil {
		return nil, code.AddSatisfactionSurveyFail
	}

	return &pb.AddCustomerSatisfactionSurveyResp{
		Id: id,
	}, nil
}
