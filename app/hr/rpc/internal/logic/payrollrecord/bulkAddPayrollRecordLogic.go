package payrollrecordlogic

import (
	"context"
	"database/sql"
	"erp/app/hr/rpc/internal/model"
	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/pb"
	"erp/common/util"
	"time"

	"erp/app/hr/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type BulkAddPayrollRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBulkAddPayrollRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BulkAddPayrollRecordLogic {
	return &BulkAddPayrollRecordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------payrollRecord-----------------------
func (l *BulkAddPayrollRecordLogic) BulkAddPayrollRecord(in *pb.BulkAddPayrollRecordReq) (*pb.BulkAddPayrollRecordResp, error) {
	var payrollRecords []*model.PayrollRecord
	for _, payrollRecord := range in.PayrollRecords {
		one, err := l.svcCtx.EmployeeDetailModel.FindOne(l.ctx, payrollRecord.EmployeeId)
		if err != nil {

			return nil, code.AddPayrollFail

		}
		netSalary := one.Salary.Float64 + payrollRecord.Bonus - payrollRecord.Deductions
		id := util.GenerateSnowflake()
		payrollRecords = append(payrollRecords, &model.PayrollRecord{
			Id:           id,
			EmployeeId:   payrollRecord.EmployeeId,
			PaymentMonth: time.Unix(payrollRecord.PaymentMonth, 0),
			BaseSalary:   one.Salary,
			Bonus:        payrollRecord.Bonus,
			Deductions:   payrollRecord.Deductions,
			NetSalary:    sql.NullFloat64{Float64: netSalary, Valid: true},
			CalculatedBy: sql.NullInt64{Int64: payrollRecord.CalculatedBy, Valid: payrollRecord.CalculatedBy != 0},
			CalculatedAt: sql.NullTime{Time: time.Unix(payrollRecord.CalculatedAt, 0), Valid: payrollRecord.CalculatedAt != 0},
			Status:       1, // 审批中
			Description:  sql.NullString{String: payrollRecord.Description, Valid: true},
			PaymentAt:    sql.NullTime{Valid: false},
		})
	}
	results, err := l.svcCtx.PayrollRecordModel.BulkInsert(payrollRecords)
	if err != nil {

		return nil, code.AddPayrollFail

	}
	var successCount, failCount int64
	var items []*pb.BulkAddPayrollRecordErrItem
	for _, r := range results {
		if r.Success {
			successCount++
		} else {
			failCount++
			logx.Error("idx:", r.EmployeeId, "err:", r.Err)

			items = append(items, &pb.BulkAddPayrollRecordErrItem{
				EmployeeId: r.EmployeeId,
				Error:      r.Err.Error(),
			})
		}
	}
	return &pb.BulkAddPayrollRecordResp{
		SuccessCount: successCount,
		ErrorCount:   failCount,
		Items:        items,
	}, nil
}
