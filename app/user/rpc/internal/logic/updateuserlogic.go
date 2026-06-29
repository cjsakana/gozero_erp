package logic

import (
	"context"
	"database/sql"
	"erp/app/user/rpc/internal/model"
	"erp/app/user/rpc/internal/svc"
	"erp/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserLogic) UpdateUser(in *pb.UpdateUserReq) (*pb.UpdateUserResp, error) {
	// 构造更新数据
	user := &model.User{
		Id:       in.Id,
		Username: in.Username,
		RealName: in.RealName,
		Phone: sql.NullString{
			String: in.Phone,
			Valid:  in.Phone != "",
		},
		Email: sql.NullString{String: in.Email, Valid: true},
		Avatar: sql.NullString{
			String: in.Avatar,
			Valid:  in.Avatar != "",
		},
	}

	// 处理离职状态
	if in.Resigned {
		user.Resigned = 1
	} else {
		user.Resigned = 0
	}

	// 设置更新人ID
	if in.UpdatedBy > 0 {
		user.UpdatedBy = sql.NullInt64{
			Int64: in.UpdatedBy,
			Valid: true,
		}
	}

	// 使用 XUpdate 进行部分更新
	err := l.svcCtx.UserModel.XUpdate(l.ctx, user)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateUserResp{}, nil
}
