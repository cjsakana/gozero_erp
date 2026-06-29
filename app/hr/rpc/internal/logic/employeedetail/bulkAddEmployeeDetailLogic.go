package employeedetaillogic

import (
	"context"
	"database/sql"
	"erp/app/hr/rpc/internal/model"
	"erp/common/encrypt"
	"erp/common/util"
	"time"

	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/pb"

	"erp/app/hr/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type BulkAddEmployeeDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBulkAddEmployeeDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BulkAddEmployeeDetailLogic {
	return &BulkAddEmployeeDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BulkAddEmployeeDetailLogic) BulkAddEmployeeDetail(in *pb.BulkAddEmployeeDetailReq) (*pb.BulkAddEmployeeDetailResp, error) {
	var employeeDetails []*model.EmployeeDetail
	for _, v := range in.EmployeeDetails {
		birthDate, _ := util.ExtractBirthdayFromID18(v.IdCard)
		idCard, _ := encrypt.EncIDCard(v.IdCard)
		account, _ := encrypt.EncAccount(v.Account)

		employeeDetails = append(employeeDetails, &model.EmployeeDetail{
			Id:           v.Id,
			EmployeeNo:   v.EmployeeNo,
			Name:         v.Name,
			IdCard:       idCard,
			Account:      sql.NullString{String: account, Valid: true},
			Gender:       v.Gender,
			BirthDate:    birthDate,
			DepartmentId: v.DepartmentId,
			PositionId:   v.PositionId,
			Salary:       sql.NullFloat64{Float64: v.Salary, Valid: true},
			HireDate:     time.Unix(v.HireDate, 0),
			LeaveDate:    sql.NullTime{Valid: false},
		})
	}

	results, err := l.svcCtx.EmployeeDetailModel.BulkInsert(employeeDetails)
	if err != nil {

		return nil, code.AddEmployeeFail

	}
	var successCount, failCount int64
	var items []*pb.BulkAddEmployeeDetailErrItem
	for _, r := range results {
		if r.Success {
			successCount++
		} else {
			failCount++
			logx.Error("idx:", r.Index, "err:", r.Err)
			items = append(items, &pb.BulkAddEmployeeDetailErrItem{
				Index: int64(r.Index),
				Error: r.Err.Error(),
			})
		}
	}

	return &pb.BulkAddEmployeeDetailResp{
		SuccessCount: successCount,
		ErrorCount:   failCount,
		Items:        items,
	}, nil
}
