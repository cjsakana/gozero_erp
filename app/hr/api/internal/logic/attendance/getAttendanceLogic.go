package attendance

import (
	"context"
	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"
	"erp/app/hr/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAttendanceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAttendanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAttendanceLogic {
	return &GetAttendanceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAttendanceLogic) GetAttendance(req *types.GetAttendanceRequest) (resp *types.GetAttendanceResponse, err error) {
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	
	// 查询考勤记录
	recordById, err := l.svcCtx.HrRPC.AttendanceRecordZrpcClient.GetAttendanceRecordById(l.ctx, &pb.GetAttendanceRecordByIdReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	
	// 查询员工信息以获取工号和姓名
	employeeDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb.GetEmployeeDetailByIdReq{
		Id: recordById.AttendanceRecord.EmployeeId,
	})
	if err != nil {
		return nil, err
	}
	
	// 构造返回数据
	resp = &types.GetAttendanceResponse{
		AttendanceRecord: types.AttendanceRecord{
			Id:            util.Int64ToString(recordById.AttendanceRecord.Id),
			EmployeeId:    util.Int64ToString(recordById.AttendanceRecord.EmployeeId),
			EmployeeNo:    employeeDetail.EmployeeNonSensitiveDetail.EmployeeNo,
			EmployeeName:  employeeDetail.EmployeeNonSensitiveDetail.Name,
			Date:          recordById.AttendanceRecord.Date,
			ClockIn:       recordById.AttendanceRecord.ClockIn,
			ClockOut:      recordById.AttendanceRecord.ClockOut,
			IsAmMissing:   recordById.AttendanceRecord.IsAmMissing,
			IsLate:        recordById.AttendanceRecord.IsLate,
			IsPmMissing:   recordById.AttendanceRecord.IsPmMissing,
			IsEarlyLeave:  recordById.AttendanceRecord.IsEarlyLeave,
			WorkHours:     recordById.AttendanceRecord.WorkHours,
			OvertimeHours: recordById.AttendanceRecord.OvertimeHours,
			Remark:        recordById.AttendanceRecord.Remark,
		},
	}

	return
}
