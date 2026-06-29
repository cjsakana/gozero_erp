package resignedapplicationlogic

import (
	"context"
	"erp/app/hr/rpc/internal/types"
	"time"

	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/pb"

	"erp/app/hr/rpc/internal/code"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchResignedApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchResignedApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchResignedApplicationLogic {
	return &SearchResignedApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchResignedApplicationLogic) SearchResignedApplication(in *pb.SearchResignedApplicationReq) (*pb.SearchResignedApplicationResp, error) {
	var startLeaveDate, endLeaveDate time.Time
	if in.StartLeaveDate > 0 {
		startLeaveDate = time.Unix(in.StartLeaveDate, 0)
	}
	if in.EndLeaveDate > 0 {
		endLeaveDate = time.Unix(in.EndLeaveDate, 0)
	}

	applications, total, err := l.svcCtx.ResignedApplicationModel.Search(l.ctx, &types.SearchResignedApplicationParams{
		SearchCom: types.SearchCom{
			Page:  in.Page,
			Limit: in.Limit,
		},
		EmployeeId:     in.EmployeeId,
		StartLeaveDate: startLeaveDate,
		EndLeaveDate:   endLeaveDate,
		Status:         in.Status,
		ApproverId:     in.ApproverId, // 使用审批人ID（新版主键）
	})
	if err != nil {
		return nil, code.GetResignFail
	}

	var pbApplications []*pb.ResignedApplication
	for _, application := range applications {
		pbApplications = append(pbApplications, &pb.ResignedApplication{
			Id:            application.Id,
			EmployeeId:    application.EmployeeId,
			Reason:        application.Reason,
			LeaveDate:     application.LeaveDate.Unix(),
			Evidence:      application.Evidence.String,
			Status:        application.Status,
			ApproverId:    application.ApproverId.Int64, // 使用审批人ID
			ApproveTime:   application.ApproveTime.Time.Unix(),
			ApproveRemark: application.ApproveRemark.String,
			CreatedAt:     application.CreatedAt.Unix(),
		})
	}
	return &pb.SearchResignedApplicationResp{
		Total:               total,
		ResignedApplication: pbApplications,
	}, nil
}
