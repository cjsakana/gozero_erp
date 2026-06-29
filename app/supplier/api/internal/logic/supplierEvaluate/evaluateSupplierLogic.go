package supplierEvaluate

import (
	"context"
	"erp/app/supplier/rpc/pb"
	"erp/common/util"
	"erp/common/xtypes"

	"erp/app/supplier/api/internal/svc"
	"erp/app/supplier/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EvaluateSupplierLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEvaluateSupplierLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EvaluateSupplierLogic {
	return &EvaluateSupplierLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EvaluateSupplierLogic) EvaluateSupplier(req *types.EvaluateSupplierReq) (resp *types.EvaluateSupplierResp, err error) {
	userId, err := util.GetInt64FromCtx(l.ctx, xtypes.UserIdKey)
	if err != nil {
		return nil, err
	}

	supplierId, err := util.StringToInt64(req.SupplierId)
	if err != nil {
		return nil, err
	}

	overallScore := req.QualityScore*0.5 + req.DeliveryScore*0.3 + req.ServiceScore*0.2
	ret, err := l.svcCtx.SupplierRPC.AddSupplierEvaluation(l.ctx, &pb.AddSupplierEvaluationReq{
		SupplierId:    supplierId,
		QualityScore:  req.QualityScore,
		DeliveryScore: req.DeliveryScore,
		ServiceScore:  req.ServiceScore,
		OverallScore:  overallScore,
		EvaluatorId:   userId,
		Remark:        req.Remark,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.EvaluateSupplierResp{
		EvaluationId: util.Int64ToString(ret.Id),
	}

	return
}
