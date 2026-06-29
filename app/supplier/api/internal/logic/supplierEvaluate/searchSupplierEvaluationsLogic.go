package supplierEvaluate

import (
	"context"
	"erp/app/hr/rpc/client/employeedetail"
	pb2 "erp/app/hr/rpc/pb"
	"erp/app/supplier/api/internal/svc"
	"erp/app/supplier/api/internal/types"
	"erp/app/supplier/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchSupplierEvaluationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchSupplierEvaluationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchSupplierEvaluationsLogic {
	return &SearchSupplierEvaluationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchSupplierEvaluationsLogic) SearchSupplierEvaluations(req *types.SearchSupplierEvaluationsReq) (resp *types.SearchSupplierEvaluationsResp, err error) {
	supplierId, err := util.StringToInt64(req.SupplierId)
	if err != nil {
		return nil, err
	}
	evaluatorId, err := util.StringToInt64(req.EvaluatorId)
	if err != nil {
		return nil, err
	}

	ret, err := l.svcCtx.SupplierRPC.SearchSupplierEvaluation(l.ctx, &pb.SearchSupplierEvaluationReq{
		Page:        req.Page,
		Limit:       req.Limit,
		SupplierId:  supplierId,
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
		EvaluatorId: evaluatorId,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.SearchSupplierEvaluationsResp{
		Total: ret.Total,
	}

	employeeMap := make(map[int64]*employeedetail.EmployeeNonSensitiveDetail)

	for _, evaluation := range ret.SupplierEvaluation {
		// 查询评估人信息
		if _, ok := employeeMap[evaluation.EvaluatorId]; !ok {
			employeeDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb2.GetEmployeeDetailByIdReq{
				Id: evaluation.EvaluatorId,
			})
			if err != nil {
				logx.Errorf("查询员工信息失败: employeeId=%d, err=%v", evaluation.EvaluatorId, err)
				continue
			}
			employeeMap[evaluation.EvaluatorId] = employeeDetail.EmployeeNonSensitiveDetail
		}

		resp.Evaluations = append(resp.Evaluations, &types.SupplierEvaluation{
			Id:            util.Int64ToString(evaluation.Id),
			SupplierId:    util.Int64ToString(evaluation.SupplierId),
			QualityScore:  evaluation.QualityScore,
			DeliveryScore: evaluation.DeliveryScore,
			ServiceScore:  evaluation.ServiceScore,
			OverallScore:  evaluation.OverallScore,
			EvaluatorId:   util.Int64ToString(evaluation.EvaluatorId),
			EvaluatorNo:   employeeMap[evaluation.EvaluatorId].EmployeeNo,
			EvaluatorName: employeeMap[evaluation.EvaluatorId].Name,
			CreatedAt:     evaluation.CreatedAt,
			Remark:        evaluation.Remark,
		})
	}

	return
}
