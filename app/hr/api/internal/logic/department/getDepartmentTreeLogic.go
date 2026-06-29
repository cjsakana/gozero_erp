package department

import (
	"context"
	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"
	"erp/app/hr/rpc/pb"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDepartmentTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDepartmentTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDepartmentTreeLogic {
	return &GetDepartmentTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDepartmentTreeLogic) GetDepartmentTree(req *types.GetDepartmentTreeRequest) (resp *types.GetDepartmentTreeResponse, err error) {
	ret, err := l.svcCtx.HrRPC.DepartmentZrpcClient.SearchDepartment(l.ctx, &pb.SearchDepartmentReq{
		Limit: -1,
	})
	if err != nil {
		return nil, err
	}

	// 构建树形结构
	treeMap := make(map[int64]*types.DepartmentTreeNode)
	var roots []*types.DepartmentTreeNode

	// 第一遍：创建所有节点
	for _, dept := range ret.Department {
		treeMap[dept.Id] = &types.DepartmentTreeNode{
			Department: types.Department{
				Id:          util.Int64ToString(dept.Id),
				Name:        dept.Name,
				ParentId:    util.Int64ToString(dept.ParentId),
				Code:        dept.Code,
				ManagerId:   util.Int64ToString(dept.ManagerId),
				ManagerNo:   dept.ManagerNo,
				ManagerName: dept.ManagerName,
				CreatedAt:   dept.CreatedAt,
				UpdatedAt:   dept.UpdatedAt,
			},
			Children: make([]*types.DepartmentTreeNode, 0),
		}
	}

	// 第二遍：建立父子关系
	for _, dept := range ret.Department {
		node := treeMap[dept.Id]

		if dept.ParentId == 0 || dept.ParentId == -1 {
			// 根节点
			roots = append(roots, node)
		} else {
			// 子节点，添加到父节点的Children中
			if parent, exists := treeMap[dept.ParentId]; exists {
				parent.Children = append(parent.Children, node)
			} else {
				// 父节点不存在，也作为根节点处理
				roots = append(roots, node)
			}
		}
	}

	// 处理指定根节点
	rootId, err := util.StringToInt64(req.RootId)
	if err != nil {
		return nil, err
	}
	if rootId > 0 {
		if rootNode, exists := treeMap[rootId]; exists {
			roots = []*types.DepartmentTreeNode{rootNode}
		} else {
			roots = []*types.DepartmentTreeNode{}
		}
	}

	return &types.GetDepartmentTreeResponse{
		Tree: roots,
	}, nil
}
