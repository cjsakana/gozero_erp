package attendance

import (
	"context"
	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"
	"erp/app/hr/rpc/client/employeedetail"
	"erp/app/hr/rpc/pb"
	"erp/common/util"
	"erp/common/xcode"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchAttendanceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchAttendanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchAttendanceLogic {
	return &SearchAttendanceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchAttendanceLogic) SearchAttendance(req *types.SearchAttendanceRequest) (resp *types.SearchAttendanceResponse, err error) {
	employeeId, err := util.StringToInt64(req.EmployeeId)
	if err != nil {
		return nil, xcode.New(1000, "参数有误")
	}

	searchReq := &pb.SearchAttendanceRecordReq{
		Page:       req.Page,
		Limit:      req.Limit,
		EmployeeId: employeeId,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		Remark:     req.Remark,
	}

	// 处理布尔筛选字段
	if req.IsLate {
		isLate := true
		searchReq.IsLate = &isLate
	}
	if req.IsEarlyLeave {
		isEarlyLeave := true
		searchReq.IsEarlyLeave = &isEarlyLeave
	}
	if req.IsAmMissing {
		isAmMissing := true
		searchReq.IsAmMissing = &isAmMissing
	}
	if req.IsPmMissing {
		isPmMissing := true
		searchReq.IsPmMissing = &isPmMissing
	}

	// 调用RPC查询
	ret, err := l.svcCtx.HrRPC.AttendanceRecordZrpcClient.SearchAttendanceRecord(l.ctx, searchReq)
	if err != nil {
		return nil, err
	}

	// 构造返回数据
	resp = &types.SearchAttendanceResponse{
		Total: ret.Total,
		List:  make([]*types.AttendanceRecord, 0, len(ret.AttendanceRecord)),
	}

	// 批量查询员工信息（为了填充员工姓名）
	employeeMap := make(map[int64]*employeedetail.EmployeeNonSensitiveDetail)
	for _, record := range ret.AttendanceRecord {
		if _, ok := employeeMap[record.EmployeeId]; !ok {
			employeeDetail, err := l.svcCtx.HrRPC.EmployeeDetailZrpcClient.GetEmployeeDetailById(l.ctx, &pb.GetEmployeeDetailByIdReq{
				Id: record.EmployeeId,
			})
			if err != nil {
				// 如果查询失败，跳过该记录或使用默认值
				logx.Errorf("查询员工信息失败: employeeId=%d, err=%v", record.EmployeeId, err)
				continue
			}
			employeeMap[record.EmployeeId] = employeeDetail.EmployeeNonSensitiveDetail
		}

		employee := employeeMap[record.EmployeeId]
		resp.List = append(resp.List, &types.AttendanceRecord{
			Id:            util.Int64ToString(record.Id),
			EmployeeId:    util.Int64ToString(record.EmployeeId),
			EmployeeNo:    employee.EmployeeNo,
			EmployeeName:  employee.Name,
			Date:          record.Date,
			ClockIn:       record.ClockIn,
			ClockOut:      record.ClockOut,
			IsAmMissing:   record.IsAmMissing,
			IsLate:        record.IsLate,
			IsPmMissing:   record.IsPmMissing,
			IsEarlyLeave:  record.IsEarlyLeave,
			WorkHours:     record.WorkHours,
			OvertimeHours: record.OvertimeHours,
			Remark:        record.Remark,
		})
	}
	return
}
