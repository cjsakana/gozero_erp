package permission

import (
	"context"
	"erp/app/auth/api/internal/svc"
	"erp/app/auth/api/internal/types"
	"erp/app/auth/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllPermissionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAllPermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllPermissionsLogic {
	return &GetAllPermissionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAllPermissionsLogic) GetAllPermissions() (resp *types.GetAllPermissionsResponse, err error) {
	ret, err := l.svcCtx.PermissionRPC.GetAllPermissions(l.ctx, &pb.GetAllPermissionsReq{})
	if err != nil {
		return nil, err
	}

	// RPC 层已经返回树形结构，直接转换类型即可
	permissionTree := l.convertPermissionTree(ret.Permissions)

	// 构造返回结果
	resp = &types.GetAllPermissionsResponse{
		Total:       ret.Total,
		Permissions: permissionTree,
	}
	return resp, nil
}

// 将 pb.PermissionItem 树形结构转换为 types.PermissionItem 树形结构
func (l *GetAllPermissionsLogic) convertPermissionTree(items []*pb.PermissionItem) []types.PermissionItem {
	if len(items) == 0 {
		return []types.PermissionItem{}
	}

	result := make([]types.PermissionItem, 0, len(items))
	for _, item := range items {
		result = append(result, l.convertPermissionItem(item))
	}
	return result
}

// 递归转换单个权限项（包含子节点）
func (l *GetAllPermissionsLogic) convertPermissionItem(item *pb.PermissionItem) types.PermissionItem {
	converted := types.PermissionItem{
		Id:          util.Int64ToString(item.Id),
		ParentId:    util.Int64ToString(item.ParentId),
		Code:        item.Code,
		Description: item.Description,
		Url:         item.Url,
		Method:      item.Method,
	}

	// 递归转换子节点
	if len(item.Child) > 0 {
		converted.Child = make([]types.PermissionItem, 0, len(item.Child))
		for _, child := range item.Child {
			converted.Child = append(converted.Child, l.convertPermissionItem(child))
		}
	} else {
		converted.Child = []types.PermissionItem{}
	}

	return converted
}
