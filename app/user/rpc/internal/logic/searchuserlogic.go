package logic

import (
	"context"
	"erp/app/user/rpc/internal/types"

	"erp/app/user/rpc/internal/svc"
	"erp/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchUserLogic {
	return &SearchUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchUserLogic) SearchUser(in *pb.SearchUserReq) (*pb.SearchUserResp, error) {
	// 构建搜索参数
	params := &types.SearchUserParams{
		Page:       in.Page,
		Limit:      in.Limit,
		EmployeeId: in.EmployeeId,
		Username:   in.Username,
		RealName:   in.RealName,
		Phone:      in.Phone,
		Email:      in.Email,
	}

	// 处理 resigned
	//if in.Resigned != nil {
	//	params.Resigned = in.Resigned
	//}
	// 暂不接受查询离职员工信息
	flag := false
	params.Resigned = &flag

	// 调用 Model 层查询
	users, total, err := l.svcCtx.UserModel.SearchUsers(l.ctx, params)
	if err != nil {
		return nil, err
	}

	// 构造返回数据
	resp := &pb.SearchUserResp{
		Total: total,
		Users: make([]*pb.User, 0, len(users)),
	}

	for _, user := range users {
		resp.Users = append(resp.Users, &pb.User{
			Id:         user.Id,
			EmployeeId: user.EmployeeId,
			EmployeeNo: user.EmployeeNo,
			Username:   user.Username,
			RealName:   user.RealName,
			Phone:      user.Phone.String,
			Email:      user.Email.String,
			Avatar:     user.Avatar.String,
			Resigned:   user.Resigned != 0,
			CreatedBy:  user.CreatedBy.Int64,
			CreatedAt:  user.CreatedAt.Unix(),
			UpdatedAt:  user.UpdatedAt.Unix(),
			UpdatedBy:  user.UpdatedBy.Int64,
		})
	}

	return resp, nil
}
