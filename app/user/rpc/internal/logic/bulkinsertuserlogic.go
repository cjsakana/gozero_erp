package logic

import (
	"context"
	"database/sql"
	"erp/app/user/rpc/internal/model"
	"erp/common/encrypt"
	"erp/common/util"
	"time"

	"erp/app/user/rpc/internal/svc"
	"erp/app/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type BulkInsertUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBulkInsertUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BulkInsertUserLogic {
	return &BulkInsertUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BulkInsertUserLogic) BulkInsertUser(in *pb.BulkInsertUserReq) (*pb.BulkInsertUserResp, error) {
	// 批量创建用户
	var users []*model.User
	var successIds []int64
	var failedIds []int64

	for _, item := range in.Users {
		// 生成用户ID
		id := util.GenerateSnowflake()

		// 密码加密
		passwordHash, err := encrypt.HashPassword(item.Password)
		if err != nil {
			logx.Errorf("密码加密失败: employeeId=%d, err=%v", item.EmployeeId, err)
			failedIds = append(failedIds, item.EmployeeId)
			continue
		}

		// 构造用户数据
		user := &model.User{
			Id:           id,
			EmployeeId:   item.EmployeeId,
			EmployeeNo:   item.EmployeeNo,
			Username:     item.Username,
			RealName:     item.RealName,
			PasswordHash: passwordHash,
			Phone: sql.NullString{
				String: item.Phone,
				Valid:  item.Phone != "",
			},
			Email:    sql.NullString{String: item.Email, Valid: item.Email != ""},
			Resigned: 0, // 默认未离职
			CreatedBy: sql.NullInt64{
				Int64: item.CreatedBy,
				Valid: item.CreatedBy > 0,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		users = append(users, user)
		successIds = append(successIds, id)
	}

	// 批量插入数据库
	if len(users) > 0 {
		err := l.svcCtx.UserModel.BulkInsert(l.ctx, users)
		if err != nil {
			return nil, err
		}
	}

	return &pb.BulkInsertUserResp{}, nil
}
