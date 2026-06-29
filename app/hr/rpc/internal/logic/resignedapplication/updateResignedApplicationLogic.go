package resignedapplicationlogic

import (
	"context"
	"database/sql"
	"erp/app/hr/rpc/internal/code"
	"erp/app/hr/rpc/internal/model"
	"time"

	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateResignedApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateResignedApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateResignedApplicationLogic {
	return &UpdateResignedApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateResignedApplicationLogic) UpdateResignedApplication(in *pb.UpdateResignedApplicationReq) (*pb.UpdateResignedApplicationResp, error) {
	err := l.svcCtx.ResignedApplicationModel.XUpdate(l.ctx, &model.ResignedApplication{
		Id:            in.Id,
		Reason:        in.Reason,
		LeaveDate:     time.Unix(in.LeaveDate, 0),
		Evidence:      sql.NullString{String: in.Evidence, Valid: in.Evidence != ""},
		Status:        in.Status,
		ApproverId:    sql.NullInt64{Int64: in.ApproverId, Valid: in.ApproverId != 0},
		ApproveTime:   sql.NullTime{Time: time.Unix(in.ApproveTime, 0), Valid: in.ApproveTime != 0},
		ApproveRemark: sql.NullString{String: in.ApproveRemark, Valid: in.ApproveRemark != ""},
	})
	if err != nil {
		return nil, code.UpdateResignFail
	}

	return &pb.UpdateResignedApplicationResp{}, nil
}
