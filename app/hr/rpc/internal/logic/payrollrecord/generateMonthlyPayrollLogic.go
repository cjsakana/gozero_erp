package payrollrecordlogic

import (
	"context"
	"database/sql"
	"time"

	"erp/app/hr/rpc/internal/model"
	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/internal/types"
	"erp/app/hr/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

// 生成指定月份工资单的逻辑
// 规则：
// - 计薪周期为自然月（上月1日至上月末），固定周六日休息（不纳入工作日计算）
// - 日薪 = 月薪 / 工作日天数；时薪 = 日薪 / 8。
// - 加班工资 = 加班小时数 * 时薪 * 3。
// - 扣款基准（可调）：缺卡（半天）100元/次；迟到/早退 50元/次。
// - 请假折算（仅统计已通过）：1-年假 100%；2-病假 80%；3-事假 0%。
// - 如果某工作日无考勤记录，且无请假覆盖，按缺上下班各一次（两次缺卡）。
// - 补卡（已通过）覆盖当日对应半天的缺卡扣款。
// - 生成记录状态：2-已核算，description=YYYY-MM 工资核算。

type GenerateMonthlyPayrollLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGenerateMonthlyPayrollLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateMonthlyPayrollLogic {
	return &GenerateMonthlyPayrollLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GenerateMonthlyPayrollLogic) GenerateMonthlyPayroll(in *pb.GenerateMonthlyPayrollReq) (*pb.GenerateMonthlyPayrollResp, error) {
	// 1) 计算目标月份（默认上月）
	year, month := resolveYearMonth(in.Year, in.Month)
	loc, _ := time.LoadLocation("Asia/Shanghai")
	start := time.Date(int(year), time.Month(month), 1, 0, 0, 0, 0, loc)
	end := start.AddDate(0, 1, 0).Add(-time.Second)

	// 2) 统计该月工作日（排除周六周日）
	workingDays := calcWorkingDays(start, end)
	if workingDays == 0 {
		return &pb.GenerateMonthlyPayrollResp{}, nil
	}

	// 3) 查询在职员工（未离职）
	emps, _, err := l.svcCtx.EmployeeDetailModel.Search(l.ctx, &types.SearchEmployeeDetailParam{
		SearchCom: types.SearchCom{Page: 1, Limit: -1},
		Resigned:  0,
	})
	if err != nil {
		return nil, err
	}

	// 4) 逐个员工计算
	success, fail := int64(0), int64(0)
	failedIds := []int64{}

	for _, emp := range emps {
		base := emp.Salary.Float64
		if base <= 0 {
			// 无基础薪资则跳过
			continue
		}
		daySalary := base / float64(workingDays)
		hourSalary := daySalary / 8.0

		// 查询该员工当月考勤（使用员工ID）
		records, _, err := l.svcCtx.AttendanceRecordModel.Search(l.ctx, &types.SearchAttendanceRecordParams{
			SearchCom:  types.SearchCom{Page: 1, Limit: -1},
			EmployeeId: emp.Id,
			StartDate:  start,
			EndDate:    end,
		})
		if err != nil {
			fail++
			failedIds = append(failedIds, emp.Id)
			continue
		}

		// 查询补卡（已通过）
		repls, _, err := l.svcCtx.AttendanceReplenishModel.Search(l.ctx, &types.SearchReplenishParams{
			SearchCom:  types.SearchCom{Page: 1, Limit: -1},
			EmployeeId: emp.Id,
			Status:     2, // 已通过
		})
		if err != nil {
			fail++
			failedIds = append(failedIds, emp.Id)
			continue
		}
		replMap := buildReplenishMap(repls)

		// 查询请假（已通过）
		leaves, _, lerr := l.svcCtx.LeaveApplicationModel.Search(l.ctx, &types.SearchLeaveApplicationParams{
			SearchCom:  types.SearchCom{Page: 1, Limit: -1},
			EmployeeId: emp.Id,
			Status:     2, // 已通过
			StartTime:  start,
			EndTime:    end,
		})
		var leaveDays map[string]leaveType
		if lerr == nil {
			leaveDays = buildLeaveDays(leaves, start, end)
		}
		if leaveDays == nil {
			leaveDays = map[string]leaveType{}
		}

		// 统计考勤、加班、缺卡与迟到早退
		recordMap := make(map[string]*model.AttendanceRecord) // YYYY-MM-DD -> record
		var overtimeHours float64
		for _, r := range records {
			key := r.Date.Format("2006-01-02")
			recordMap[key] = r
			overtimeHours += r.OvertimeHours
		}

		// 逐工作日汇总扣款
		const (
			penaltyMissingHalf = 100.0 // 缺卡半天
			penaltyLateEarly   = 50.0  // 迟到/早退
		)
		var deductions float64
		for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
			if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
				continue
			}
			key := d.Format("2006-01-02")

			// 请假覆盖（按天），根据类型计付比例
			if lt, ok := leaveDays[key]; ok {
				paidRatio := leavePaidRatio(lt)
				// 未付部分计入扣款
				deductions += daySalary * (1 - paidRatio)
				continue // 请假当天不再按缺卡/迟到早退罚
			}

			r, ok := recordMap[key]
			if !ok {
				// 无记录：视为上下班各缺一次
				if !isReplenished(replMap[key], true) { // 上午
					deductions += penaltyMissingHalf
				}
				if !isReplenished(replMap[key], false) { // 下午
					deductions += penaltyMissingHalf
				}
				continue
			}

			// 迟到/早退（使用新版布尔字段）
			if r.IsLate == 1 {
				deductions += penaltyLateEarly
			}
			if r.IsEarlyLeave == 1 {
				deductions += penaltyLateEarly
			}
			// 缺卡半天（若补卡覆盖则不扣）
			if r.IsAmMissing == 1 && !isReplenished(replMap[key], true) {
				deductions += penaltyMissingHalf
			}
			if r.IsPmMissing == 1 && !isReplenished(replMap[key], false) {
				deductions += penaltyMissingHalf
			}
		}

		// 加班奖金
		bonus := overtimeHours * hourSalary * 3

		net := base + bonus - deductions

		// 生成雪花ID
		id := util.GenerateSnowflake()
		rec := &model.PayrollRecord{
			Id:           id,
			EmployeeId:   emp.Id, // 使用员工ID
			PaymentMonth: start,  // 薪资月份
			BaseSalary:   sql.NullFloat64{Float64: base, Valid: true},
			Bonus:        bonus,
			Deductions:   deductions,
			NetSalary:    sql.NullFloat64{Float64: net, Valid: true},
			CalculatedBy: sql.NullInt64{Int64: 0, Valid: false}, // system
			CalculatedAt: sql.NullTime{Time: time.Now(), Valid: true},
			Status:       2, // 已核算
			Description:  sql.NullString{String: start.Format("2006-01") + " 工资核算", Valid: true},
		}
		if _, err := l.svcCtx.PayrollRecordModel.Insert(l.ctx, rec); err != nil {
			fail++
			failedIds = append(failedIds, emp.Id)
			continue
		}
		success++
	}

	return &pb.GenerateMonthlyPayrollResp{SuccessCount: success, FailCount: fail, FailedEmployeeIds: failedIds}, nil
}

func resolveYearMonth(year, month int64) (int64, int64) {
	if year > 0 && month > 0 {
		return year, month
	}
	now := time.Now()
	pm := now.AddDate(0, -1, 0)
	return int64(pm.Year()), int64(pm.Month())
}

func calcWorkingDays(start, end time.Time) int {
	cnt := 0
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
			continue
		}
		cnt++
	}
	return cnt
}

type leaveType int

// 1-年假 2-病假 3-事假
func leavePaidRatio(t leaveType) float64 {
	switch t {
	case 1:
		return 1.0
	case 2:
		return 0.8
	case 3:
		return 0.0
	default:
		return 0.0
	}
}

// buildLeaveDays 从请假申请构建"天->请假类型"的映射（按自然日覆盖）
func buildLeaveDays(list []*model.LeaveApplication, start, end time.Time) map[string]leaveType {
	m := map[string]leaveType{}
	for _, v := range list {
		st := v.StartTime
		ed := v.EndTime
		if st.After(end) || ed.Before(start) {
			continue
		}
		if st.Before(start) {
			st = start
		}
		if ed.After(end) {
			ed = end
		}
		for d := st; !d.After(ed); d = d.AddDate(0, 0, 1) {
			key := d.Format("2006-01-02")
			m[key] = leaveType(v.Type)
		}
	}
	return m
}

// replInfo 同一天补卡合并信息
type replInfo struct {
	morning   bool
	afternoon bool
}

// buildReplenishMap 构建补卡映射
func buildReplenishMap(list []*model.AttendanceReplenish) map[string]replInfo {
	m := map[string]replInfo{}
	for _, r := range list {
		key := r.OriginalDate.Format("2006-01-02")
		cur := m[key]
		switch r.ReplenishType {
		case 1: // 上班
			cur.morning = true
		case 2: // 下班
			cur.afternoon = true
		case 3: // 全天
			cur.morning, cur.afternoon = true, true
		}
		m[key] = cur
	}
	return m
}

// isReplenished 判断是否已补卡
func isReplenished(info replInfo, morning bool) bool {
	if morning {
		return info.morning
	}
	return info.afternoon
}
