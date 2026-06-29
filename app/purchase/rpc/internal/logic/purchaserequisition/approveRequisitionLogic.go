package purchaserequisitionlogic

import (
	"context"
	"erp/app/purchase/rpc/internal/types"

	"erp/app/purchase/rpc/internal/code"
	"erp/app/purchase/rpc/internal/svc"
	"erp/app/purchase/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApproveRequisitionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApproveRequisitionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveRequisitionLogic {
	return &ApproveRequisitionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 审批采购申请
func (l *ApproveRequisitionLogic) ApproveRequisition(in *pb.ApproveRequisitionReq) (*pb.ApproveRequisitionResp, error) {
	err := l.svcCtx.PurchaseRequisitionModel.Approve(l.ctx,
		&types.ApproveRequisitionParam{
			Id:            in.RequisitionId,
			ApproveTime:   in.ApproveTime,
			ApproveRemark: in.ApproveRemark,
			TargetStatus:  in.TargetStatus,
		},
	)
	if err != nil {
		return nil, code.ApproveRequisitionFail
	}
	return &pb.ApproveRequisitionResp{}, nil
}
