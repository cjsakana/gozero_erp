package user

import (
	"context"
	"erp/app/user/api/internal/code"
	"erp/app/user/rpc/pb"
	"erp/common/util"
	"erp/common/xtypes"
	"fmt"
	"net/http"

	"erp/app/user/api/internal/svc"
	"erp/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadAvatarLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadAvatarLogic {
	return &UploadAvatarLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

const maxFileSize = 5 * 1024 * 1024 // 5MB

func (l *UploadAvatarLogic) UploadAvatar(req *http.Request) (resp *types.EmptyResponse, err error) {
	userid, err := util.GetInt64FromCtx(l.ctx, xtypes.UserIdKey)
	if err != nil {
		return nil, err
	}

	// formData 获取图片
	_ = req.ParseMultipartForm(maxFileSize)
	file, handler, err := req.FormFile("avatar")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fmt.Println("1111111", handler.Filename)

	url, err := l.svcCtx.UploadClient.UploadFile(l.ctx, file, handler.Filename)
	if err != nil {
		return nil, code.PutBucketErr
	}
	fmt.Println("2222222")
	// 更新DB
	_, err = l.svcCtx.UserRPC.UpdateUser(l.ctx, &pb.UpdateUserReq{
		Id:     userid,
		Avatar: url,
	})
	fmt.Println("8958565")
	if err != nil {
		fmt.Println("3333333333", err)
		return nil, err
	}

	return
}
