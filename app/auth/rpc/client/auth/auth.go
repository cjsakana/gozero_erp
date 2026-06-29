package auth

import (
	"erp/app/auth/rpc/client/permission"
	"erp/app/auth/rpc/client/role"
	"erp/app/auth/rpc/client/rolepermission"
	"erp/app/auth/rpc/client/userrole"
	"github.com/zeromicro/go-zero/zrpc"
)

type (
	AuthZrpcClient struct {
		permission.PermissionZrpcClient
		role.RoleZrpcClient
		rolepermission.RolePermissionZrpcClient
		userrole.UserRoleZrpcClient
	}
)

func NewAuthZrpcClient(cli zrpc.Client) AuthZrpcClient {
	return AuthZrpcClient{
		PermissionZrpcClient:     permission.NewPermissionZrpcClient(cli),
		RoleZrpcClient:           role.NewRoleZrpcClient(cli),
		RolePermissionZrpcClient: rolepermission.NewRolePermissionZrpcClient(cli),
		UserRoleZrpcClient:       userrole.NewUserRoleZrpcClient(cli),
	}
}
