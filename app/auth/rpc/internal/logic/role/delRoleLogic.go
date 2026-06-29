package rolelogic

import (
	"context"
	"fmt"

	"erp/app/auth/rpc/internal/code"
	"erp/app/auth/rpc/internal/svc"
	"erp/app/auth/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelRoleLogic {
	return &DelRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelRoleLogic) DelRole(in *pb.DelRoleReq) (*pb.DelRoleResp, error) {
	// 先查询角色信息，用于后续缓存清理
	data, err := l.svcCtx.RoleModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, code.DeleteRoleFail
	}

	// 使用事务删除角色和角色权限
	rpIds, err := l.svcCtx.RoleModel.XDeleteTX(l.ctx, in.Id)
	if err != nil {
		return nil, code.DeleteRoleFail
	}

	// 事务成功后，清理缓存
	cacheErpAuthRoleIdPrefix := "cache:erpAuth:role:id:"
	cacheErpAuthRoleCodePrefix := "cache:erpAuth:role:code:"
	erpAuthRoleCodeKey := fmt.Sprintf("%s%v", cacheErpAuthRoleCodePrefix, data.Code)
	erpAuthRoleIdKey := fmt.Sprintf("%s%v", cacheErpAuthRoleIdPrefix, in.Id)
	l.svcCtx.BizRedis.Del(erpAuthRoleCodeKey, erpAuthRoleIdKey)

	// 清理角色权限缓存
	for _, id := range rpIds {
		cacheErpAuthRolePermissionIdPrefix := "cache:erpAuth:rolePermission:id:"
		erpAuthRolePermissionIdKey := fmt.Sprintf("%s%v", cacheErpAuthRolePermissionIdPrefix, id)
		l.svcCtx.BizRedis.Del(erpAuthRolePermissionIdKey)
	}

	return &pb.DelRoleResp{}, nil
}
