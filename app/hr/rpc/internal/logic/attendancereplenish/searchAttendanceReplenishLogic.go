package attendancereplenishlogic

import (
	"context"
	"erp/app/hr/rpc/internal/svc"
	types2 "erp/app/hr/rpc/internal/types"
	"erp/app/hr/rpc/pb"

	"erp/app/hr/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type SearchAttendanceReplenishLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchAttendanceReplenishLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchAttendanceReplenishLogic {
	return &SearchAttendanceReplenishLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchAttendanceReplenishLogic) SearchAttendanceReplenish(in *pb.SearchAttendanceReplenishReq) (*pb.SearchAttendanceReplenishResp, error) {
	attendanceReplenishes, total, err := l.svcCtx.AttendanceReplenishModel.Search(l.ctx, &types2.SearchReplenishParams{
		SearchCom: types2.SearchCom{
			Page:  in.Page,
			Limit: in.Limit,
		},
		EmployeeId:    in.EmployeeId, // 使用员工ID
		ReplenishType: in.ReplenishType,
		Reason:        in.Reason,
		Status:        in.Status,
		ApproverId:    in.ApproverId, // 使用审批人ID
	})
	if err != nil {
		return nil, code.SearchAttendanceFail
	}

	var pbAttendanceReplenishes []*pb.AttendanceReplenish
	for _, attendanceReplenish := range attendanceReplenishes {
		pbAttendanceReplenishes = append(pbAttendanceReplenishes, &pb.AttendanceReplenish{
			Id:            attendanceReplenish.Id,
			EmployeeId:    attendanceReplenish.EmployeeId,
			OriginalDate:  attendanceReplenish.OriginalDate.Unix(),
			ApplyTime:     attendanceReplenish.ApplyTime.Unix(),
			ReplenishType: attendanceReplenish.ReplenishType,
			ReplenishTime: attendanceReplenish.ReplenishTime.Time.Unix(),
			Reason:        attendanceReplenish.Reason,
			Evidence:      attendanceReplenish.Evidence.String,
			Status:        attendanceReplenish.Status,
			ApproverId:    attendanceReplenish.ApproverId.Int64,
			ApproveTime:   attendanceReplenish.ApproveTime.Time.Unix(),
			ApproveRemark: attendanceReplenish.ApproveRemark.String,
		})
	}

	return &pb.SearchAttendanceReplenishResp{
		Total:               total,
		AttendanceReplenish: pbAttendanceReplenishes,
	}, nil
}
