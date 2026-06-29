package employeedetaillogic

import (
	"context"
	"erp/app/hr/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/stores/sqlc"

	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmployeeDetailByEmployeeNoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEmployeeDetailByEmployeeNoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmployeeDetailByEmployeeNoLogic {
	return &GetEmployeeDetailByEmployeeNoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetEmployeeDetailByEmployeeNo 根据员工no获取员工详情
func (l *GetEmployeeDetailByEmployeeNoLogic) GetEmployeeDetailByEmployeeNo(in *pb.GetEmployeeDetailByEmployeeNoReq) (*pb.GetEmployeeDetailByEmployeeNoResp, error) {
	one, err := l.svcCtx.EmployeeDetailModel.FindOneByEmployeeNo(l.ctx, in.EmployeeNo)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, code.EmployeeNotFound
		}
		return nil, err
	}

	return &pb.GetEmployeeDetailByEmployeeNoResp{
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
