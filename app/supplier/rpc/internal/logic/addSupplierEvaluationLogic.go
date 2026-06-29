package logic

import (
	"context"
	"database/sql"
	"erp/app/supplier/rpc/internal/code"
	"erp/app/supplier/rpc/internal/model"
	"erp/app/supplier/rpc/internal/svc"
	"erp/app/supplier/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddSupplierEvaluationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddSupplierEvaluationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddSupplierEvaluationLogic {
	return &AddSupplierEvaluationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------supplierEvaluation-----------------------
func (l *AddSupplierEvaluationLogic) AddSupplierEvaluation(in *pb.AddSupplierEvaluationReq) (*pb.AddSupplierEvaluationResp, error) {
	ret, err := l.svcCtx.EvaluationModel.Insert(l.ctx, &model.SupplierEvaluation{
		SupplierId:    in.SupplierId,
		QualityScore:  in.QualityScore,
		DeliveryScore: in.DeliveryScore,
		ServiceScore:  in.ServiceScore,
		OverallScore:  in.OverallScore,
		EvaluatorId:   in.EvaluatorId,
		Remark:        sql.NullString{String: in.Remark, Valid: true},
	})
	if err != nil {
		return nil, code.AddEvaluationFail
	}

	id, _ := ret.LastInsertId()
	return &pb.AddSupplierEvaluationResp{
		Id: id,
	}, nil
}
