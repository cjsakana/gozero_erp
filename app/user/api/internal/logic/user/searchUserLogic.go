package user

import (
	"context"
	"erp/app/user/api/internal/svc"
	"erp/app/user/api/internal/types"
	"erp/app/user/rpc/user"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchUserLogic {
	return &SearchUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchUserLogic) SearchUser(req *types.SearchUserReq) (resp *types.SearchUserResp, err error) {
	ret, err := l.svcCtx.UserRPC.SearchUser(l.ctx, &user.SearchUserReq{
		Page:     req.Page,
		Limit:    req.Limit,
		Username: req.Username,
		RealName: req.RealName,
		Phone:    req.Phone,
		Email:    req.Email,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.SearchUserResp{
		Total: ret.Total,
	}
	var users []types.User

	for _, u := range ret.Users {
		users = append(users, types.User{
			Id:         util.Int64ToString(u.Id),
			EmployeeId: util.Int64ToString(u.EmployeeId),
			EmployeeNo: u.EmployeeNo,
			Username:   u.Username,
			RealName:   u.RealName,
			Phone:      u.Phone,
			Email:      u.Email,
			Avatar:     u.Avatar,
			Resigned:   u.Resigned,
			CreatedBy:  util.Int64ToString(u.CreatedBy),
			CreatedAt:  u.CreatedAt,
			UpdatedAt:  u.UpdatedAt,
			UpdatedBy:  util.Int64ToString(u.UpdatedBy),
		})
	}
	resp.Users = users

	return
}
