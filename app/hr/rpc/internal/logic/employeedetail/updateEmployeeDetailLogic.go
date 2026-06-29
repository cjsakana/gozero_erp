package employeedetaillogic

import (
	"context"
	"database/sql"
	"erp/app/hr/rpc/internal/model"
	"time"

	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/pb"

	"erp/app/hr/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateEmployeeDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateEmployeeDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateEmployeeDetailLogic {
	return &UpdateEmployeeDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateEmployeeDetailLogic) UpdateEmployeeDetail(in *pb.UpdateEmployeeDetailReq) (*pb.UpdateEmployeeDetailResp, error) {
	err := l.svcCtx.EmployeeDetailModel.XUpdate(l.ctx, &model.EmployeeDetail{
		Id:           in.Id,
		Account:      sql.NullString{String: in.Account, Valid: true},
		DepartmentId: in.DepartmentId,
		PositionId:   in.PositionId,
		Salary:       sql.NullFloat64{Float64: in.Salary, Valid: true},
		HireDate:     time.Unix(in.HireDate, 0),
		LeaveDate:    sql.NullTime{Time: time.Unix(int64(in.LeaveDate), 0), Valid: true},
		Name:         in.Name,
	})
	if err != nil {

		return nil, code.UpdateEmployeeFail

	}

	return &pb.UpdateEmployeeDetailResp{}, nil
}
