package image

import (
	"context"
	"erp/app/image/api/internal/svc"
	"erp/app/image/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
)

type UploadImageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadImageLogic {
	return &UploadImageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

const maxFileSize = 100 * 1024 * 1024 // 100MB

func (l *UploadImageLogic) UploadImage(r *http.Request, req *types.UploadImageReq) (resp *types.UploadImageResp, err error) {
	// formData 获取多张图片，保持顺序
	_ = r.ParseMultipartForm(maxFileSize)
	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		return nil, http.ErrMissingFile
	}
	files := r.MultipartForm.File["images"]
	if len(files) == 0 {
		return nil, http.ErrMissingFile
	}

	items := make([]*types.UploadImageItem, 0, len(files))
	for i, fh := range files {
		// 打开当前文件
		f, err := fh.Open()
		if err != nil {
			return nil, err
		}
		// 上传
		url, err := l.svcCtx.UploadClient.UploadFile(l.ctx, f, fh.Filename)
		_ = f.Close()
		if err != nil {
			return nil, err
		}
		// 追加结果，Index 为顺序下标
		items = append(items, &types.UploadImageItem{
			Url:   url,
			Index: int64(i),
		})
	}

	return &types.UploadImageResp{Item: items}, nil
}
