package logic

import (
	"context"
	types2 "erp/app/customer/rpc/internal/types"

	"erp/app/customer/rpc/internal/code"
	"erp/app/customer/rpc/internal/svc"
	"erp/app/customer/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchCustomerSatisfactionSurveyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchCustomerSatisfactionSurveyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchCustomerSatisfactionSurveyLogic {
	return &SearchCustomerSatisfactionSurveyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchCustomerSatisfactionSurveyLogic) SearchCustomerSatisfactionSurvey(in *pb.SearchCustomerSatisfactionSurveyReq) (*pb.SearchCustomerSatisfactionSurveyResp, error) {
	evaluations, total, err := l.svcCtx.SatisfactionModel.Search(l.ctx, &types2.SearchSatisfaction{
		SearchCom: types2.SearchCom{
			Page:  in.Page,
			Limit: in.Limit,
		},
		CustomerId:  in.CustomerId,
		QualityMin:  in.QualityMin,
		QualityMax:  in.QualityMax,
		QualityOp:   in.QualityOp,
		DeliveryMin: in.DeliveryMin,
		DeliveryMax: in.DeliveryMax,
		DeliveryOp:  in.DeliveryOp,
		ServiceMin:  in.ServiceMin,
		ServiceMax:  in.ServiceMax,
		ServiceOp:   in.ServiceOp,
		OverallMin:  in.OverallMin,
		OverallMax:  in.OverallMax,
		OverallOp:   in.OverallOp,
		StartData:   in.StartDate,
		EndData:     in.EndDate,
	})
	if err != nil {
		return nil, code.GetSatisfactionSurveyFail
	}

	list := make([]*pb.CustomerSatisfactionSurvey, total)
	for i, one := range evaluations {
		list[i] = &pb.CustomerSatisfactionSurvey{
			Id:            one.Id,
			CustomerId:    one.CustomerId,
			QualityScore:  one.QualityScore,
			DeliveryScore: one.DeliveryScore,
			ServiceScore:  one.ServiceScore,
			OverallScore:  one.OverallScore,
			Remark:        one.Remark.String,
			CreatedAt:     one.CreatedAt.Unix(),
			CreatedBy:     one.CreatedBy,
		}
	}

	return &pb.SearchCustomerSatisfactionSurveyResp{
		CustomerSatisfactionSurvey: list,
		Total:                      total,
	}, nil
}
