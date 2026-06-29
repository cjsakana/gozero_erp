package purchaseOrder

import (
	"context"
	"erp/app/purchase/api/internal/code"
	"erp/app/purchase/rpc/pb"
	"erp/common/util"
	"net/http"

	"erp/app/purchase/api/internal/svc"
	"erp/app/purchase/api/internal/types"

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
func (l *UploadContractFileLogic) UploadContractFile(r *http.Request, req *types.UploadContractFileRequest) (resp *types.UploadContractFileResponse, err error) {

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
	id, err := util.StringToInt64(req.Id)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.PurchaseRPC.UpdateOrder(l.ctx, &pb.UpdateOrderReq{
		Id:          id,
		ContractUrl: url,
	})
	if err != nil {
		return nil, err
	}

	return
}
