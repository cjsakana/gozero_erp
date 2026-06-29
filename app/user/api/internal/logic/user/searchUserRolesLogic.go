package user

import (
	"context"
	"erp/app/user/api/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type SearchUserRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchUserRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchUserRolesLogic {
	return &SearchUserRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

//func (l *SearchUserRolesLogic) SearchUserRoles(req *types.SearchUserRolesRequest) (resp *types.SearchUserRolesResponse, err error) {
//	// 设置默认分页参数
//	page := req.Page
//	limit := req.Limit
//	if page <= 0 {
//		page = 1
//	}
//	if limit <= 0 {
//		limit = 10
//	}
//
//	// 判断是否提供了用户筛选条件
//	hasUserFilter := req.EmployeeNo != "" || req.RealName != ""
//
//	var userMap map[int64]*pb.User
//	var userIdSet map[int64]bool
//
//	// 如果提供了用户筛选条件，先筛选用户
//	if hasUserFilter {
//		ret1, err := l.svcCtx.UserRPC.SearchUser(l.ctx, &pb.SearchUserReq{
//			Limit:      -1,
//			EmployeeNo: req.EmployeeNo,
//			RealName:   req.RealName,
//		})
//		if err != nil {
//			return nil, err
//		}
//
//		// 如果没有符合条件的用户，直接返回空结果
//		if len(ret1.Users) == 0 {
//			return &types.SearchUserRolesResponse{
//				Total: 0,
//				Item:  []*types.SearchUserRolesItem{},
//			}, nil
//		}
//
//		userMap = make(map[int64]*pb.User)
//		userIdSet = make(map[int64]bool)
//		for _, user := range ret1.Users {
//			userMap[user.Id] = user
//			userIdSet[user.Id] = true
//		}
//	}
//
//	// 如果提供了用户筛选条件，需要查询所有符合条件的用户角色关系（不分页），然后在内存中过滤和分页
//	// 如果没有提供用户筛选条件，可以直接使用分页查询
//	if hasUserFilter {
//		// 查询所有符合条件的用户角色关系（不分页）
//		ret2, err := l.svcCtx.UserRoleRPC.SearchUserRole(l.ctx, &pb2.SearchUserRoleReq{
//			Page:   1,
//			Limit:  -1, // -1 表示查询全部
//			RoleId: req.RoleId,
//		})
//		if err != nil {
//			return nil, err
//		}
//
//		// 在内存中过滤：只保留符合条件的用户
//		var filteredItems []*types.SearchUserRolesItem
//		for _, item := range ret2.UserRole {
//			// 只保留在 userIdSet 中的用户
//			if !userIdSet[item.UserId] {
//				continue
//			}
//			u := userMap[item.UserId]
//			filteredItems = append(filteredItems, &types.SearchUserRolesItem{
//				Id:         item.Id,
//				EmployeeNo: u.EmployeeNo,
//				RealName:   u.RealName,
//				RoleId:     item.RoleId,
//			})
//		}
//
//		// 在内存中做分页
//		total := int64(len(filteredItems))
//		start := (page - 1) * limit
//		end := start + limit
//		if start > total {
//			start = total
//		}
//		if end > total {
//			end = total
//		}
//
//		var items []*types.SearchUserRolesItem
//		if start < end {
//			items = filteredItems[start:end]
//		}
//
//		return &types.SearchUserRolesResponse{
//			Total: total,
//			Item:  items,
//		}, nil
//	} else {
//		// 没有提供用户筛选条件，直接使用分页查询
//		ret2, err := l.svcCtx.UserRoleRPC.SearchUserRole(l.ctx, &pb2.SearchUserRoleReq{
//			Page:   page,
//			Limit:  limit,
//			RoleId: req.RoleId,
//		})
//		if err != nil {
//			return nil, err
//		}
//
//		// 需要查询用户信息来填充 EmployeeNo 和 RealName
//		if len(ret2.UserRole) > 0 {
//			// 收集所有用户ID
//			userIds := make([]int64, 0, len(ret2.UserRole))
//			userIdMap := make(map[int64]bool)
//			for _, item := range ret2.UserRole {
//				if !userIdMap[item.UserId] {
//					userIds = append(userIds, item.UserId)
//					userIdMap[item.UserId] = true
//				}
//			}
//
//			// 批量查询用户信息
//			userMap = make(map[int64]*pb.User)
//			for _, userId := range userIds {
//				userResp, err := l.svcCtx.UserRPC.GetUserById(l.ctx, &pb.GetUserByIdReq{Id: userId})
//				if err == nil && userResp != nil {
//					userMap[userId] = userResp.User
//				}
//			}
//		}
//
//		var items []*types.SearchUserRolesItem
//		for _, item := range ret2.UserRole {
//			u, ok := userMap[item.UserId]
//			if !ok {
//				// 如果用户不存在，跳过或使用默认值
//				continue
//			}
//			items = append(items, &types.SearchUserRolesItem{
//				Id:         item.Id,
//				UserId:     u.Id,
//				EmployeeNo: u.EmployeeNo,
//				RealName:   u.RealName,
//				RoleId:     item.RoleId,
//			})
//		}
//
//		return &types.SearchUserRolesResponse{
//			Total: ret2.Total,
//			Item:  items,
//		}, nil
//	}
//}
