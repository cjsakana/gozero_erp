package leaveapplicationlogic

import (
	"context"
	"erp/app/hr/rpc/internal/code"
	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/pb"
	"github.com/zeromicro/go-zero/core/stores/sqlc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLeaveApplicationByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLeaveApplicationByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLeaveApplicationByIdLogic {
	return &GetLeaveApplicationByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetLeaveApplicationByIdLogic) GetLeaveApplicationById(in *pb.GetLeaveApplicationByIdReq) (*pb.GetLeaveApplicationByIdResp, error) {
	one, err := l.svcCtx.LeaveApplicationModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, code.LeaveNotFound
		}
		return nil, code.LeaveNotFound
	}

	return &pb.GetLeaveApplicationByIdResp{
		LeaveApplication: &pb.LeaveApplication{
			Id:            one.Id,
			EmployeeId:    one.EmployeeId,
			Type:          one.Type,
			StartTime:     one.StartTime.Unix(),
			EndTime:       one.EndTime.Unix(),
			Duration:      one.Duration,
			Reason:        one.Reason,
			Evidence:      one.Evidence.String,
			Status:        one.Status,
			ApproverId:    one.ApproverId.Int64,
			ApproveTime:   one.ApproveTime.Time.Unix(),
			ApproveRemark: one.ApproveRemark.String,
			CreatedAt:     one.CreatedAt.Unix(),
		},
	}, nil
}
