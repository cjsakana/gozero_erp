package purchaserequisitionlogic

import (
	"context"

	"erp/app/purchase/rpc/internal/svc"
	"erp/app/purchase/rpc/internal/types"
	"erp/app/purchase/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRequisitionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateRequisitionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRequisitionLogic {
	return &UpdateRequisitionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新采购申请（动态SQL拼接）
func (l *UpdateRequisitionLogic) UpdateRequisition(in *pb.UpdateRequisitionReq) (*pb.UpdateRequisitionResp, error) {
	// 构建更新参数
	param := &types.UpdateRequisitionParam{
		Id: in.Id,
	}

	// 只有非零值才设置指针
	if in.DepartmentId != 0 {
		param.DepartmentId = &in.DepartmentId
	}
	if in.ApplicantId != 0 {
		param.ApplicantId = &in.ApplicantId
	}
	if in.RequestDate != 0 {
		param.RequestDate = &in.RequestDate
	}
	if in.TotalAmount != 0 {
		param.TotalAmount = &in.TotalAmount
	}
	if in.Status != 0 {
		param.Status = &in.Status
	}
	if in.ApproverId != 0 {
		param.ApproverId = &in.ApproverId
	}
	if in.ApproveTime != 0 {
		param.ApproveTime = &in.ApproveTime
	}
	if in.ApproveRemark != "" {
		param.ApproveRemark = &in.ApproveRemark
	}

	// 调用model层更新方法
	err := l.svcCtx.PurchaseRequisitionModel.UpdateRequisition(l.ctx, param)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateRequisitionResp{}, nil
}
