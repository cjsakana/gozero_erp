package employeedetaillogic

import (
	"context"
	"erp/app/hr/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/stores/sqlc"

	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmployeeDetailByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEmployeeDetailByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmployeeDetailByIdLogic {
	return &GetEmployeeDetailByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetEmployeeDetailByIdLogic) GetEmployeeDetailById(in *pb.GetEmployeeDetailByIdReq) (*pb.GetEmployeeDetailByIdResp, error) {
	one, err := l.svcCtx.EmployeeDetailModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, code.EmployeeNotFound
		}
		return nil, err
	}

	return &pb.GetEmployeeDetailByIdResp{
		EmployeeNonSensitiveDetail: &pb.EmployeeNonSensitiveDetail{
			Id:           one.Id,
			EmployeeNo:   one.EmployeeNo,
			Name:         one.Name,
			Gender:       one.Gender,
			BirthDate:    one.BirthDate.Unix(),
			DepartmentId: one.DepartmentId,
			PositionId:   one.PositionId,
			Salary:       one.Salary.Float64,
			HireDate:     one.HireDate.Unix(),
			LeaveDate:    one.LeaveDate.Time.Unix(),
		},
	}, nil

}
