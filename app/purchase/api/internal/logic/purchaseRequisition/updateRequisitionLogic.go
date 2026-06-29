package purchaseRequisition

import (
	"context"
	"erp/common/util"

	"erp/app/purchase/api/internal/svc"
	"erp/app/purchase/api/internal/types"
	"erp/app/purchase/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRequisitionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateRequisitionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRequisitionLogic {
	return &UpdateRequisitionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRequisitionLogic) UpdateRequisition(req *types.UpdateRequisitionReq) (resp *types.UpdateRequisitionResp, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	departmentId, err := util.StringToInt64(req.DepartmentId)
	if err != nil {
		return nil, err
	}
	applicantId, err := util.StringToInt64(req.ApplicantId)
	if err != nil {
		return nil, err
	}
	approverId, err := util.StringToInt64(req.ApproverId)
	if err != nil {
		return nil, err
	}

	// 调用RPC服务
	_, err = l.svcCtx.PurchaseRPC.UpdateRequisition(l.ctx, &pb.UpdateRequisitionReq{
		Id:            id,
		DepartmentId:  departmentId,
		ApplicantId:   applicantId,
		RequestDate:   req.RequestDate,
		TotalAmount:   req.TotalAmount,
		Status:        req.Status,
		ApproverId:    approverId,
		ApproveTime:   req.ApproveTime,
		ApproveRemark: req.ApproveRemark,
	})
	if err != nil {
		return nil, err
	}

	return &types.UpdateRequisitionResp{}, nil
}
