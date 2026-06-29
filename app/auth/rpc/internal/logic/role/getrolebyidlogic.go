package rolelogic

import (
	"context"

	"erp/app/auth/rpc/internal/svc"
	"erp/app/auth/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoleByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRoleByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleByIdLogic {
	return &GetRoleByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRoleByIdLogic) GetRoleById(in *pb.GetRoleByIdReq) (*pb.GetRoleByIdResp, error) {
	one, err := l.svcCtx.RoleModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetRoleByIdResp{
		Role: &pb.Role{
			Id:          one.Id,
			Code:        one.Code,
			Name:        one.Name,
			Description: one.Description.String,
		},
	}, nil
}
