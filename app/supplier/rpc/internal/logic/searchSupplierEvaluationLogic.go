package logic

import (
	"context"
	"erp/app/supplier/rpc/internal/code"
	"erp/app/supplier/rpc/internal/svc"
	types2 "erp/app/supplier/rpc/internal/types"
	"erp/app/supplier/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchSupplierEvaluationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchSupplierEvaluationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchSupplierEvaluationLogic {
	return &SearchSupplierEvaluationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchSupplierEvaluationLogic) SearchSupplierEvaluation(in *pb.SearchSupplierEvaluationReq) (*pb.SearchSupplierEvaluationResp, error) {
	evaluations, total, err := l.svcCtx.EvaluationModel.Search(l.ctx, &types2.SearchSupplierEvaluation{
		SearchCom: types2.SearchCom{
			Page:  in.Page,
			Limit: in.Limit,
		},
		SupplierId:  in.SupplierId,
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
		EvaluatorId: in.EvaluatorId,
	})
	if err != nil {
		return nil, code.SearchEvaluationFail
	}

	list := make([]*pb.SupplierEvaluation, 0)
	for _, one := range evaluations {
		list = append(list, &pb.SupplierEvaluation{
			Id:            one.Id,
			SupplierId:    one.SupplierId,
			QualityScore:  one.QualityScore,
			DeliveryScore: one.DeliveryScore,
			ServiceScore:  one.ServiceScore,
			OverallScore:  one.OverallScore,
			EvaluatorId:   one.EvaluatorId,
			Remark:        one.Remark.String,
			CreatedAt:     one.CreatedAt.Unix(),
		})
	}

	return &pb.SearchSupplierEvaluationResp{
		SupplierEvaluation: list,
		Total:              total,
	}, nil
}
