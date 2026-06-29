package employeedetaillogic

import (
	"context"
	types2 "erp/app/hr/rpc/internal/types"
	"time"

	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/pb"

	"erp/app/hr/rpc/internal/code"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchEmployeeDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchEmployeeDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchEmployeeDetailLogic {
	return &SearchEmployeeDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchEmployeeDetailLogic) SearchEmployeeDetail(in *pb.SearchEmployeeDetailReq) (*pb.SearchEmployeeDetailResp, error) {
	var hireDate time.Time
	if in.HireDate > 0 {
		hireDate = time.Unix(in.HireDate, 0)
	}

	records, total, err := l.svcCtx.EmployeeDetailModel.Search(l.ctx, &types2.SearchEmployeeDetailParam{
		SearchCom: types2.SearchCom{
			Page:  in.Page,
			Limit: in.Limit,
		},
		Gender:       in.Gender,
		DepartmentId: in.DepartmentId,
		PositionId:   in.PositionId,
		Salary:       in.Salary,
		HireDate:     hireDate,
		Name:         in.Name,
		Resigned:     in.Resigned,
	})
	if err != nil {

		return nil, code.GetEmployeeFail

	}

	var res []*pb.EmployeeNonSensitiveDetail
	for _, record := range records {
		res = append(res, &pb.EmployeeNonSensitiveDetail{
			Id:           record.Id,
			EmployeeNo:   record.EmployeeNo,
			Name:         record.Name,
			Gender:       record.Gender,
			BirthDate:    record.BirthDate.Unix(),
			DepartmentId: record.DepartmentId,
			PositionId:   record.PositionId,
			Salary:       record.Salary.Float64,
			HireDate:     record.HireDate.Unix(),
			LeaveDate:    record.LeaveDate.Time.Unix(),
		})
	}
	return &pb.SearchEmployeeDetailResp{
		Total:                      total,
		EmployeeNonSensitiveDetail: res,
	}, nil
}
