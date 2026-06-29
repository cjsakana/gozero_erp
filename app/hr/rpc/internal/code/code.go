package code

import "erp/common/xcode"

var (
	// 员工相关错误 110101-110110
	EmployeeNotFound      = xcode.New(110101, "员工不存在")
	EmployeeAlreadyExists = xcode.New(110102, "员工已存在")
	EmployeeNoDuplicate   = xcode.New(110103, "员工编号重复")
	AddEmployeeFail       = xcode.New(110104, "添加员工失败")
	UpdateEmployeeFail    = xcode.New(110105, "更新员工失败")
	DeleteEmployeeFail    = xcode.New(110106, "删除员工失败")
	GetEmployeeFail       = xcode.New(110107, "获取员工信息失败")
	ImportEmployeeFail    = xcode.New(110108, "导入员工失败")
	AdjustSalaryFail      = xcode.New(110109, "调整薪资失败")

	// 部门相关错误 110111-110120
	DepartmentNotFound      = xcode.New(110111, "部门不存在")
	DepartmentAlreadyExists = xcode.New(110112, "部门已存在")
	DepartmentCodeDuplicate = xcode.New(110113, "部门编码重复")
	AddDepartmentFail       = xcode.New(110114, "添加部门失败")
	UpdateDepartmentFail    = xcode.New(110115, "更新部门失败")
	DeleteDepartmentFail    = xcode.New(110116, "删除部门失败")
	DepartmentInUse         = xcode.New(110117, "部门正在使用中，无法删除")
	GetDepartmentFail       = xcode.New(110118, "获取部门信息失败")
	ImportDepartmentFail    = xcode.New(110119, "导入部门失败")

	// 职位相关错误 110121-110130
	PositionNotFound        = xcode.New(110121, "职位不存在")
	PositionAlreadyExists   = xcode.New(110122, "职位已存在")
	PositionCodeDuplicate   = xcode.New(110123, "职位编码重复")
	AddPositionFail         = xcode.New(110124, "添加职位失败")
	UpdatePositionFail      = xcode.New(110125, "更新职位失败")
	DeletePositionFail      = xcode.New(110126, "删除职位失败")
	PositionInUse           = xcode.New(110127, "职位正在使用中，无法删除")
	GetPositionFail         = xcode.New(110128, "获取职位信息失败")
	BatchAssignPositionFail = xcode.New(110129, "批量分配职位失败")

	// 考勤相关错误 110131-110140
	AttendanceNotFound   = xcode.New(110131, "考勤记录不存在")
	ClockFail            = xcode.New(110132, "打卡失败")
	GetAttendanceFail    = xcode.New(110133, "获取考勤记录失败")
	SearchAttendanceFail = xcode.New(110134, "搜索考勤记录失败")
	DuplicateClock       = xcode.New(110135, "重复打卡")
	ClockTimeInvalid     = xcode.New(110136, "打卡时间无效")

	// 考勤补卡相关错误 110141-110150
	AttendanceReplenishNotFound = xcode.New(110141, "考勤补卡申请不存在")
	SubmitReplenishFail         = xcode.New(110142, "提交考勤补卡申请失败")
	ApproveReplenishFail        = xcode.New(110143, "审批考勤补卡申请失败")
	GetReplenishFail            = xcode.New(110144, "获取考勤补卡申请失败")
	ReplenishStatusInvalid      = xcode.New(110145, "补卡申请状态无效")
	ReplenishAlreadyApproved    = xcode.New(110146, "补卡申请已审批")

	// 请假相关错误 110151-110160
	LeaveNotFound        = xcode.New(110151, "请假申请不存在")
	SubmitLeaveFail      = xcode.New(110152, "提交请假申请失败")
	ApproveLeaveFail     = xcode.New(110153, "审批请假申请失败")
	GetLeaveFail         = xcode.New(110154, "获取请假申请失败")
	LeaveStatusInvalid   = xcode.New(110155, "请假申请状态无效")
	LeaveAlreadyApproved = xcode.New(110156, "请假申请已审批")
	LeaveDaysInvalid     = xcode.New(110157, "请假天数无效")

	// 工资单相关错误 110161-110170
	PayrollNotFound        = xcode.New(110161, "工资单不存在")
	AddPayrollFail         = xcode.New(110162, "添加工资单失败")
	ApprovePayrollFail     = xcode.New(110163, "审批工资单失败")
	ExecutePaymentFail     = xcode.New(110164, "执行工资发放失败")
	SubmitToFinanceFail    = xcode.New(110165, "提交财务失败")
	GetPayrollFail         = xcode.New(110166, "获取工资单失败")
	SearchPayrollFail      = xcode.New(110167, "搜索工资单失败")
	PayrollStatusInvalid   = xcode.New(110168, "工资单状态无效")
	PayrollAlreadyApproved = xcode.New(110169, "工资单已审批")

	// 离职相关错误 110171-110180
	ResignedApplicationNotFound = xcode.New(110171, "离职申请不存在")
	ApplyResignFail             = xcode.New(110172, "申请离职失败")
	ReviewResignFail            = xcode.New(110173, "审核离职申请失败")
	UpdateResignFail            = xcode.New(110174, "更新离职申请失败")
	GetResignFail               = xcode.New(110175, "获取离职申请失败")
	ResignStatusInvalid         = xcode.New(110176, "离职申请状态无效")
	ResignAlreadyApproved       = xcode.New(110177, "离职申请已审批")
)
