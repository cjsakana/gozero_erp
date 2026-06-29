package productCategory

import (
	"context"
	"erp/app/product/rpc/pb"
	"erp/common/util"
	"sort"

	"erp/app/product/api/internal/svc"
	"erp/app/product/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCategoryTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCategoryTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCategoryTreeLogic {
	return &GetCategoryTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCategoryTreeLogic) GetCategoryTree(req *types.GetCategoryTreeRequest) (resp *types.GetCategoryTreeResponse, err error) {
	ret, err := l.svcCtx.ProductRPC.SearchProductCategory(l.ctx, &pb.SearchProductCategoryReq{
		Limit: -1,
	})
	if err != nil {
		return nil, err
	}

	// 构建树形结构（使用指针）
	treeMap := make(map[string]*types.ProductCategory)
	var roots []*types.ProductCategory

	// 第一遍：创建所有节点
	for _, one := range ret.ProductCategory {

		treeMap[util.Int64ToString(one.CategoryId)] = &types.ProductCategory{
			CategoryId:   util.Int64ToString(one.CategoryId),
			CategoryName: one.CategoryName,
			ParentId:     util.Int64ToString(one.ParentId),
			Children:     []*types.ProductCategory{},
		}
	}

	// 第二遍：建立父子关系
	for _, category := range treeMap {
		if category.ParentId == "0" || category.ParentId == "-1" {
			// 根节点
			roots = append(roots, category)
		} else {
			// 子节点，添加到父节点的Children中
			if parent, exists := treeMap[category.ParentId]; exists {
				parent.Children = append(parent.Children, category)
			} else {
				// 父节点不存在，也作为根节点处理
				roots = append(roots, category)
			}
		}
	}

	// 处理指定根节点
	parentId, err := util.StringToInt64(req.ParentId)
	if err != nil {
		return nil, err
	}
	if parentId > 0 {
		if rootNode, exists := treeMap[req.ParentId]; exists {
			roots = []*types.ProductCategory{rootNode}
		} else {
			roots = []*types.ProductCategory{}
		}
	}

	// 排序
	sortedRoots := l.sortTree(roots)

	return &types.GetCategoryTreeResponse{
		Categories: sortedRoots,
	}, nil

}

// 排序函数（使用指针）
func (l *GetCategoryTreeLogic) sortTree(categories []*types.ProductCategory) []*types.ProductCategory {
	sort.Slice(categories, func(i, j int) bool {
		return categories[i].CategoryName < categories[j].CategoryName
	})

	for _, category := range categories {
		if len(category.Children) > 0 {
			category.Children = l.sortTree(category.Children)
		}
	}

	return categories
}
