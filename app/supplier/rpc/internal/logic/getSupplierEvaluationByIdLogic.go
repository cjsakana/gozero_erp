package logic

import (
	"context"

	"erp/app/supplier/rpc/internal/code"
	"erp/app/supplier/rpc/internal/svc"
	"erp/app/supplier/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
)

type GetSupplierEvaluationByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSupplierEvaluationByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSupplierEvaluationByIdLogic {
	return &GetSupplierEvaluationByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// rpc UpdateSupplierEvaluation(UpdateSupplierEvaluationReq) returns (UpdateSupplierEvaluationResp);
func (l *GetSupplierEvaluationByIdLogic) GetSupplierEvaluationById(in *pb.GetSupplierEvaluationByIdReq) (*pb.GetSupplierEvaluationByIdResp, error) {
	one, err := l.svcCtx.EvaluationModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, code.SupplierEvaluationNotFound
		}
		return nil, code.GetEvaluationFail
	}

	return &pb.GetSupplierEvaluationByIdResp{
		SupplierEvaluation: &pb.SupplierEvaluation{
			Id:            one.Id,
			SupplierId:    one.SupplierId,
			QualityScore:  one.QualityScore,
			DeliveryScore: one.DeliveryScore,
			ServiceScore:  one.ServiceScore,
			OverallScore:  one.OverallScore,
			EvaluatorId:   one.EvaluatorId,
			Remark:        one.Remark.String,
			CreatedAt:     one.CreatedAt.Unix(),
		},
	}, nil
}
