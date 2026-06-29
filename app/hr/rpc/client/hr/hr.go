package hr

import (
	"erp/app/hr/rpc/client/attendancerecord"
	"erp/app/hr/rpc/client/attendancereplenish"
	"erp/app/hr/rpc/client/department"
	"erp/app/hr/rpc/client/employeedetail"
	"erp/app/hr/rpc/client/leaveapplication"
	"erp/app/hr/rpc/client/payrollrecord"
	"erp/app/hr/rpc/client/position"
	"erp/app/hr/rpc/client/resignedapplication"
	"github.com/zeromicro/go-zero/zrpc"
)

type HrZrpcClient struct {
	attendancerecord.AttendanceRecordZrpcClient
	attendancereplenish.AttendanceReplenishZrpcClient
	department.DepartmentZrpcClient
	employeedetail.EmployeeDetailZrpcClient
	leaveapplication.LeaveApplicationZrpcClient
	payrollrecord.PayrollRecordZrpcClient
	position.PositionZrpcClient
	resignedapplication.ResignedApplicationZrpcClient
}

func NewHrZrpcClient(cli zrpc.Client) HrZrpcClient {
	return HrZrpcClient{
		AttendanceRecordZrpcClient:    attendancerecord.NewAttendanceRecordZrpcClient(cli),
		AttendanceReplenishZrpcClient: attendancereplenish.NewAttendanceReplenishZrpcClient(cli),
		DepartmentZrpcClient:          department.NewDepartmentZrpcClient(cli),
		EmployeeDetailZrpcClient:      employeedetail.NewEmployeeDetailZrpcClient(cli),
		LeaveApplicationZrpcClient:    leaveapplication.NewLeaveApplicationZrpcClient(cli),
		PayrollRecordZrpcClient:       payrollrecord.NewPayrollRecordZrpcClient(cli),
		PositionZrpcClient:            position.NewPositionZrpcClient(cli),
		ResignedApplicationZrpcClient: resignedapplication.NewResignedApplicationZrpcClient(cli),
	}
}
