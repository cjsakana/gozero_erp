package resignedapplicationlogic

import (
	"context"

	"erp/app/hr/rpc/internal/svc"
	"erp/app/hr/rpc/pb"

	"erp/app/hr/rpc/internal/code"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
)

type GetResignedApplicationByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetResignedApplicationByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetResignedApplicationByIdLogic {
	return &GetResignedApplicationByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetResignedApplicationByIdLogic) GetResignedApplicationById(in *pb.GetResignedApplicationByIdReq) (*pb.GetResignedApplicationByIdResp, error) {
	one, err := l.svcCtx.ResignedApplicationModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == sqlc.ErrNotFound {
			return nil, code.ResignedApplicationNotFound
		}
		return nil, code.ResignedApplicationNotFound

	}

	return &pb.GetResignedApplicationByIdResp{
		ResignedApplication: &pb.ResignedApplication{
			Id:            one.Id,
			EmployeeId:    one.EmployeeId,
			Reason:        one.Reason,
			LeaveDate:     one.LeaveDate.Unix(),
			Evidence:      one.Evidence.String,
			Status:        one.Status,
			ApproverId:    one.ApproverId.Int64,
			ApproveTime:   one.ApproveTime.Time.Unix(),
			ApproveRemark: one.ApproveRemark.String,
			CreatedAt:     one.CreatedAt.Unix(),
		},
	}, nil
}
