package logic

import (
	"context"
	"erp/app/customer/rpc/internal/code"
	"erp/app/customer/rpc/internal/svc"
	"erp/app/customer/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
)

type GetCustomerSatisfactionSurveyByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCustomerSatisfactionSurveyByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCustomerSatisfactionSurveyByIdLogic {
	return &GetCustomerSatisfactionSurveyByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCustomerSatisfactionSurveyByIdLogic) GetCustomerSatisfactionSurveyById(in *pb.GetCustomerSatisfactionSurveyByIdReq) (*pb.GetCustomerSatisfactionSurveyByIdResp, error) {
	one, err := l.svcCtx.SatisfactionModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, code.CustomerSatisfactionNotFound
		}
		return nil, code.GetSatisfactionSurveyFail
	}

	return &pb.GetCustomerSatisfactionSurveyByIdResp{
		CustomerSatisfactionSurvey: &pb.CustomerSatisfactionSurvey{
			Id:            one.Id,
			CustomerId:    one.CustomerId,
			QualityScore:  one.QualityScore,
			DeliveryScore: one.DeliveryScore,
			ServiceScore:  one.ServiceScore,
			OverallScore:  one.OverallScore,
			Remark:        one.Remark.String,
			CreatedBy:     one.CreatedBy,
			CreatedAt:     one.CreatedAt.Unix(),
		},
	}, nil
}
