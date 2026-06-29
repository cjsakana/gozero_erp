package permissionlogic

import (
	"context"
	"erp/app/auth/rpc/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"erp/app/auth/rpc/internal/svc"
	"erp/app/auth/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllPermissionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllPermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllPermissionsLogic {
	return &GetAllPermissionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------permission-----------------------
func (l *GetAllPermissionsLogic) GetAllPermissions(in *pb.GetAllPermissionsReq) (*pb.GetAllPermissionsResp, error) {
	// 1. 从数据库查询所有权限（平铺结构）
	permissions, err := l.svcCtx.PermissionModel.FindAll(l.ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "查询权限列表失败: "+err.Error())
	}

	// 2. 构建树形权限结构
	permissionTree := l.buildPermissionTree(permissions)

	// 3. 返回响应
	return &pb.GetAllPermissionsResp{
		Total:       int64(len(permissions)),
		Permissions: permissionTree,
	}, nil
}

// 构建权限树形结构
func (l *GetAllPermissionsLogic) buildPermissionTree(permissions []*model.Permission) []*pb.PermissionItem {
	// 创建权限映射表
	permissionMap := make(map[int64]*pb.PermissionItem)
	var roots []*pb.PermissionItem

	// 第一遍：创建所有权限节点
	for _, perm := range permissions {
		permissionItem := &pb.PermissionItem{
			Id:          perm.Id,
			ParentId:    perm.ParentId,
			Code:        perm.Code.String,
			Description: perm.Description.String,
			Url:         perm.Url.String,
			Method:      perm.Method.String,
			Child:       make([]*pb.PermissionItem, 0),
		}
		permissionMap[perm.Id] = permissionItem
	}

	// 第二遍：构建父子关系
	for _, perm := range permissions {
		permissionItem := permissionMap[perm.Id]
		if perm.ParentId == 0 {
			// 根节点
			roots = append(roots, permissionItem)
		} else {
			// 子节点，添加到父节点的 Child 列表中
			if parent, exists := permissionMap[perm.ParentId]; exists {
				parent.Child = append(parent.Child, permissionItem)
			}
		}
	}

	return roots
}
