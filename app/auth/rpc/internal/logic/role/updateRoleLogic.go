package rolelogic

import (
	"context"
	"database/sql"
	"erp/app/auth/rpc/internal/code"
	"erp/app/auth/rpc/internal/model"

	"erp/app/auth/rpc/internal/svc"
	"erp/app/auth/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleLogic {
	return &UpdateRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateRoleLogic) UpdateRole(in *pb.UpdateRoleReq) (*pb.UpdateRoleResp, error) {
	err := l.svcCtx.RoleModel.Update(l.ctx, &model.Role{
		Id:          in.Id,
		Code:        in.Code,
		Name:        in.Name,
		Description: sql.NullString{String: in.Description, Valid: true},
	})
	if err != nil {
		return nil, code.UpdateRoleFail
	}
	return &pb.UpdateRoleResp{}, nil
}
