package userrolelogic

import (
	"context"
	"database/sql"
	"erp/app/auth/rpc/internal/code"
	"erp/app/auth/rpc/internal/model"
	"erp/common/util"
	"fmt"

	"erp/app/auth/rpc/internal/svc"
	"erp/app/auth/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddUserRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddUserRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddUserRoleLogic {
	return &AddUserRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------userRole-----------------------
func (l *AddUserRoleLogic) AddUserRole(in *pb.AddUserRoleReq) (*pb.AddUserRoleResp, error) {
	id := util.GenerateSnowflake()
	_, err := l.svcCtx.UserRoleModel.Insert(l.ctx, &model.UserRole{
		Id:     id,
		UserId: sql.NullInt64{Int64: in.UserId, Valid: true},
		RoleId: sql.NullInt64{Int64: in.RoleId, Valid: true},
	})
	if err != nil {
		fmt.Println("err", err)
		return nil, code.AddUserRoleFail
	}

	return &pb.AddUserRoleResp{
		Id: id,
	}, nil
}
