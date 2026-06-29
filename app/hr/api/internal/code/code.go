package code

import "erp/common/xcode"

var (
	// 员工相关错误 110001-110010
	EmployeeNotFound      = xcode.New(110001, "员工不存在")
	EmployeeAlreadyExists = xcode.New(110002, "员工已存在")
	EmployeeNoDuplicate   = xcode.New(110003, "员工编号重复")
	AddEmployeeFail       = xcode.New(110004, "添加员工失败")
	UpdateEmployeeFail    = xcode.New(110005, "更新员工失败")
	DeleteEmployeeFail    = xcode.New(110006, "删除员工失败")
	GetEmployeeFail       = xcode.New(110007, "获取员工信息失败")
	ImportEmployeeFail    = xcode.New(110008, "导入员工失败")
	AdjustSalaryFail      = xcode.New(110009, "调整薪资失败")

	// 部门相关错误 110011-110020
	DepartmentNotFound      = xcode.New(110011, "部门不存在")
	DepartmentAlreadyExists = xcode.New(110012, "部门已存在")
	DepartmentCodeDuplicate = xcode.New(110013, "部门编码重复")
	AddDepartmentFail       = xcode.New(110014, "添加部门失败")
	UpdateDepartmentFail    = xcode.New(110015, "更新部门失败")
	DeleteDepartmentFail    = xcode.New(110016, "删除部门失败")
	DepartmentInUse         = xcode.New(110017, "部门正在使用中，无法删除")
	GetDepartmentFail       = xcode.New(110018, "获取部门信息失败")
	ImportDepartmentFail    = xcode.New(110019, "导入部门失败")

	// 职位相关错误 110021-110030
	PositionNotFound        = xcode.New(110021, "职位不存在")
	PositionAlreadyExists   = xcode.New(110022, "职位已存在")
	PositionCodeDuplicate   = xcode.New(110023, "职位编码重复")
	AddPositionFail         = xcode.New(110024, "添加职位失败")
	UpdatePositionFail      = xcode.New(110025, "更新职位失败")
	DeletePositionFail      = xcode.New(110026, "删除职位失败")
	PositionInUse           = xcode.New(110027, "职位正在使用中，无法删除")
	GetPositionFail         = xcode.New(110028, "获取职位信息失败")
	BatchAssignPositionFail = xcode.New(110029, "批量分配职位失败")

	// 考勤相关错误 110031-110040
	AttendanceNotFound   = xcode.New(110031, "考勤记录不存在")
	ClockFail            = xcode.New(110032, "打卡失败")
	GetAttendanceFail    = xcode.New(110033, "获取考勤记录失败")
	SearchAttendanceFail = xcode.New(110034, "搜索考勤记录失败")
	DuplicateClock       = xcode.New(110035, "重复打卡")
	ClockTimeInvalid     = xcode.New(110036, "打卡时间无效")

	// 考勤补卡相关错误 110041-110050
	AttendanceReplenishNotFound = xcode.New(110041, "考勤补卡申请不存在")
	SubmitReplenishFail         = xcode.New(110042, "提交考勤补卡申请失败")
	ApproveReplenishFail        = xcode.New(110043, "审批考勤补卡申请失败")
	GetReplenishFail            = xcode.New(110044, "获取考勤补卡申请失败")
	ReplenishStatusInvalid      = xcode.New(110045, "补卡申请状态无效")
	ReplenishAlreadyApproved    = xcode.New(110046, "补卡申请已审批")

	// 请假相关错误 110051-110060
	LeaveNotFound        = xcode.New(110051, "请假申请不存在")
	SubmitLeaveFail      = xcode.New(110052, "提交请假申请失败")
	ApproveLeaveFail     = xcode.New(110053, "审批请假申请失败")
	GetLeaveFail         = xcode.New(110054, "获取请假申请失败")
	LeaveStatusInvalid   = xcode.New(110055, "请假申请状态无效")
	LeaveAlreadyApproved = xcode.New(110056, "请假申请已审批")
	LeaveDaysInvalid     = xcode.New(110057, "请假天数无效")

	// 工资单相关错误 110061-110070
	PayrollNotFound        = xcode.New(110061, "工资单不存在")
	AddPayrollFail         = xcode.New(110062, "添加工资单失败")
	ApprovePayrollFail     = xcode.New(110063, "审批工资单失败")
	ExecutePaymentFail     = xcode.New(110064, "执行工资发放失败")
	SubmitToFinanceFail    = xcode.New(110065, "提交财务失败")
	GetPayrollFail         = xcode.New(110066, "获取工资单失败")
	SearchPayrollFail      = xcode.New(110067, "搜索工资单失败")
	PayrollStatusInvalid   = xcode.New(110068, "工资单状态无效")
	PayrollAlreadyApproved = xcode.New(110069, "工资单已审批")

	// 离职相关错误 110071-110080
	ResignedApplicationNotFound = xcode.New(110071, "离职申请不存在")
	ApplyResignFail             = xcode.New(110072, "申请离职失败")
	ReviewResignFail            = xcode.New(110073, "审核离职申请失败")
	UpdateResignFail            = xcode.New(110074, "更新离职申请失败")
	GetResignFail               = xcode.New(110075, "获取离职申请失败")
	ResignStatusInvalid         = xcode.New(110076, "离职申请状态无效")
	ResignAlreadyApproved       = xcode.New(110077, "离职申请已审批")
)
