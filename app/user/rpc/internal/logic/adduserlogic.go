package logic

import (
	"context"
	"database/sql"
	"erp/app/user/rpc/internal/model"
	"erp/app/user/rpc/internal/svc"
	"erp/app/user/rpc/pb"
	"erp/common/encrypt"
	"erp/common/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddUserLogic {
	return &AddUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// -----------------------user-----------------------
func (l *AddUserLogic) AddUser(in *pb.AddUserReq) (*pb.AddUserResp, error) {
	// 生成用户ID
	id := util.GenerateSnowflake()

	// 密码加密
	passwordHash, err := encrypt.HashPassword(in.Password)
	if err != nil {
		return nil, err
	}

	// 构造用户数据
	user := &model.User{
		Id:           id,
		EmployeeId:   in.EmployeeId,
		EmployeeNo:   in.EmployeeNo,
		Username:     in.Username,
		RealName:     in.RealName,
		PasswordHash: passwordHash,
		Phone: sql.NullString{
			String: in.Phone,
			Valid:  in.Phone != "",
		},
		Email:    sql.NullString{String: in.Email, Valid: in.Email != ""},
		Resigned: 0, // 默认未离职
		CreatedBy: sql.NullInt64{
			Int64: in.CreatedBy,
			Valid: in.CreatedBy > 0,
		},
		UpdatedBy: sql.NullInt64{
			Int64: in.CreatedBy,
			Valid: in.CreatedBy > 0,
		},
	}

	// 插入数据库
	_, err = l.svcCtx.UserModel.Insert(l.ctx, user)
	if err != nil {
		return nil, err
	}

	return &pb.AddUserResp{
		Id: user.Id,
	}, nil
}
