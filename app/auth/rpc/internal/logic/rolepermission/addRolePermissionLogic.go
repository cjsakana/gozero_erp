package rolepermissionlogic

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

type AddRolePermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddRolePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddRolePermissionLogic {
	return &AddRolePermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------rolePermission-----------------------
func (l *AddRolePermissionLogic) AddRolePermission(in *pb.AddRolePermissionReq) (*pb.AddRolePermissionResp, error) {
	id := util.GenerateSnowflake()
	_, err := l.svcCtx.RolePermissionModel.Insert(l.ctx, &model.RolePermission{
		Id:           id,
		RoleId:       sql.NullInt64{Int64: in.RoleId, Valid: true},
		PermissionId: sql.NullInt64{Int64: in.PermissionId, Valid: true},
	})
	if err != nil {
		return nil, code.AddRolePermissionFail
	}

	return &pb.AddRolePermissionResp{
		Id: id,
	}, nil
}
