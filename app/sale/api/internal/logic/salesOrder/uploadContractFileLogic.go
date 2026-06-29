package salesOrder

import (
	"context"
	"erp/app/sale/api/internal/code"
	"net/http"

	"erp/app/sale/api/internal/svc"
	"erp/app/sale/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadContractFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadContractFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadContractFileLogic {
	return &UploadContractFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

const maxFileSize = 50 * 1024 * 1024 // 5MB
func (l *UploadContractFileLogic) UploadContractFile(r *http.Request) (resp *types.UploadContractFileResponse, err error) {

	_ = r.ParseMultipartForm(maxFileSize)
	file, handler, err := r.FormFile("contract")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	url, err := l.svcCtx.UploadClient.UploadFile(l.ctx, file, handler.Filename)
	if err != nil {
		return nil, code.PutBucketErr
	}

	resp = &types.UploadContractFileResponse{
		ContractURL: url,
	}

	return
}
