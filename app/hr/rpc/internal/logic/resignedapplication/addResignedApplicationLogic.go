package resignedapplicationlogic

import (
	"context"
	"database/sql"
	"erp/app/hr/rpc/internal/model"
	"erp/common/util"
	"time"

	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/pb"

	"erp/app/hr/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddResignedApplicationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddResignedApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddResignedApplicationLogic {
	return &AddResignedApplicationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------离职申请表-----------------------
func (l *AddResignedApplicationLogic) AddResignedApplication(in *pb.AddResignedApplicationReq) (*pb.AddResignedApplicationResp, error) {
	// 生成雪花ID
	id := util.GenerateSnowflake()
	_, err := l.svcCtx.ResignedApplicationModel.Insert(l.ctx, &model.ResignedApplication{
		Id:            id,
		EmployeeId:    in.EmployeeId,
		Reason:        in.Reason,
		LeaveDate:     time.Unix(in.LeaveDate, 0),
		Evidence:      sql.NullString{String: in.Evidence, Valid: in.Evidence != ""},
		Status:        1, // 审批中
		ApproverId:    sql.NullInt64{Int64: in.ApproverId, Valid: in.ApproverId != 0},
		ApproveTime:   sql.NullTime{Valid: false},
		ApproveRemark: sql.NullString{Valid: false},
	})
	if err != nil {
		return nil, code.ApplyResignFail
	}
	return &pb.AddResignedApplicationResp{
		Id: id,
	}, nil
}
