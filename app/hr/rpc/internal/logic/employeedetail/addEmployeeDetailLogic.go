package employeedetaillogic

import (
	"context"
	"database/sql"
	"erp/app/hr/rpc/internal/code"
	"erp/app/hr/rpc/internal/model"
	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/pb"
	"erp/common/encrypt"
	"erp/common/util"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type AddEmployeeDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddEmployeeDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddEmployeeDetailLogic {
	return &AddEmployeeDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------员工信息扩展表-----------------------
func (l *AddEmployeeDetailLogic) AddEmployeeDetail(in *pb.AddEmployeeDetailReq) (*pb.AddEmployeeDetailResp, error) {
	birthDate, _ := util.ExtractBirthdayFromID18(in.EmployeeDetail.IdCard)

	idCard, _ := encrypt.EncIDCard(in.EmployeeDetail.IdCard)
	account, _ := encrypt.EncAccount(in.EmployeeDetail.Account)

	id := util.GenerateSnowflake()
	_, err := l.svcCtx.EmployeeDetailModel.Insert(l.ctx, &model.EmployeeDetail{
		Id:           in.EmployeeDetail.Id,
		EmployeeNo:   in.EmployeeDetail.EmployeeNo,
		Name:         in.EmployeeDetail.Name,
		IdCard:       idCard,
		Account:      sql.NullString{String: account, Valid: true},
		Gender:       in.EmployeeDetail.Gender,
		BirthDate:    birthDate,
		DepartmentId: in.EmployeeDetail.DepartmentId,
		PositionId:   in.EmployeeDetail.PositionId,
		Salary:       sql.NullFloat64{Float64: in.EmployeeDetail.Salary, Valid: true},
		HireDate:     time.Unix(in.EmployeeDetail.HireDate, 0),
		LeaveDate:    sql.NullTime{Valid: false},
	})
	if err != nil {

		return nil, code.AddEmployeeFail

	}

	return &pb.AddEmployeeDetailResp{
		Id: id,
	}, nil
}
