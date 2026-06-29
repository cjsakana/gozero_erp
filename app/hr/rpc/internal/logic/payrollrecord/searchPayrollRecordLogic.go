package payrollrecordlogic

import (
	"context"
	"erp/app/hr/rpc/internal/types"
	"time"

	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/pb"

	"erp/app/hr/rpc/internal/code"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchPayrollRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchPayrollRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchPayrollRecordLogic {
	return &SearchPayrollRecordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchPayrollRecordLogic) SearchPayrollRecord(in *pb.SearchPayrollRecordReq) (*pb.SearchPayrollRecordResp, error) {
	var startCalculatedDate, endCalculatedDate, startPaymentDate, endPaymentDate, paymentMonth time.Time
	if in.StartCalculatedDate > 0 {
		startCalculatedDate = time.Unix(in.StartCalculatedDate, 0)
	}
	if in.EndCalculatedDate > 0 {
		endCalculatedDate = time.Unix(in.EndCalculatedDate, 0)
	}
	if in.StartPaymentDate > 0 {
		startPaymentDate = time.Unix(in.StartPaymentDate, 0)
	}
	if in.EndPaymentDate > 0 {
		endPaymentDate = time.Unix(in.EndPaymentDate, 0)
	}
	if in.PaymentMonth != 0 {
		paymentMonth = time.Unix(in.PaymentMonth, 0)
	}

	records, total, err := l.svcCtx.PayrollRecordModel.Search(l.ctx, &types.SearchPayrollRecordParams{
		SearchCom: types.SearchCom{
			Page:  in.Page,
			Limit: in.Limit,
		},
		EmployeeId:          in.EmployeeId, // 使用员工ID
		Status:              in.Status,
		Description:         in.Description,
		CalculatedBy:        in.CalculatedBy, // 使用核算人ID
		StartCalculatedDate: startCalculatedDate,
		EndCalculatedDate:   endCalculatedDate,
		StartPaymentDate:    startPaymentDate,
		EndPaymentDate:      endPaymentDate,
		PaymentMonth:        paymentMonth,
	})
	if err != nil {
		return nil, code.SearchPayrollFail
	}

	var pbPayrollRecords []*pb.PayrollRecord
	for _, record := range records {
		pbPayrollRecords = append(pbPayrollRecords, &pb.PayrollRecord{
			Id:           record.Id,
			EmployeeId:   record.EmployeeId, // 使用员工ID
			PaymentMonth: record.PaymentMonth.Unix(),
			BaseSalary:   record.BaseSalary.Float64,
			Bonus:        record.Bonus,
			Deductions:   record.Deductions,
			NetSalary:    record.NetSalary.Float64,
			CalculatedBy: record.CalculatedBy.Int64, // 使用核算人ID
			CalculatedAt: record.CalculatedAt.Time.Unix(),
			Status:       record.Status,
			Description:  record.Description.String,
			PaymentAt:    record.PaymentAt.Time.Unix(),
			CreatedAt:    record.CreatedAt.Unix(),
		})
	}

	return &pb.SearchPayrollRecordResp{
		Total:         total,
		PayrollRecord: pbPayrollRecords,
	}, nil
}
