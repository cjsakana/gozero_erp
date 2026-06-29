package purchaseRequisition

import (
	"context"
	"erp/app/purchase/api/internal/svc"
	"erp/app/purchase/api/internal/types"
	"erp/app/purchase/rpc/pb"
	"erp/common/util"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApproveRequisitionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApproveRequisitionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveRequisitionLogic {
	return &ApproveRequisitionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApproveRequisitionLogic) ApproveRequisition(req *types.ApproveRequisitionReq) (resp *types.ApproveRequisitionResp, err error) {
	requisitionId, err := util.StringToInt64(req.RequisitionId)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.PurchaseRPC.ApproveRequisition(l.ctx, &pb.ApproveRequisitionReq{
		RequisitionId: requisitionId,
		ApproveTime:   time.Now().Unix(),
		ApproveRemark: req.ApproveRemark,
		TargetStatus:  req.TargetStatus,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.ApproveRequisitionResp{}
	return
}
