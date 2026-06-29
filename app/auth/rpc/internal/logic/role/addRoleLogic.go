package rolelogic

import (
	"context"
	"database/sql"
	"erp/app/auth/rpc/internal/code"
	"erp/app/auth/rpc/internal/model"
	"erp/common/util"

	"erp/app/auth/rpc/internal/svc"
	"erp/app/auth/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddRoleLogic {
	return &AddRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------role-----------------------
func (l *AddRoleLogic) AddRole(in *pb.AddRoleReq) (*pb.AddRoleResp, error) {
	id := util.GenerateSnowflake()
	_, err := l.svcCtx.RoleModel.Insert(l.ctx, &model.Role{
		Id:          id,
		Code:        in.Code,
		Name:        in.Name,
		Description: sql.NullString{String: in.Description, Valid: true},
	})
	if err != nil {
		return nil, code.AddRoleFail
	}

	return &pb.AddRoleResp{
		Id: id,
	}, nil
}
