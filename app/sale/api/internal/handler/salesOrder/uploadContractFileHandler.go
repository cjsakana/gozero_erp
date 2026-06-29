package salesOrder

import (
	"net/http"

	"erp/app/sale/api/internal/logic/salesOrder"
	"erp/app/sale/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UploadContractFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := salesOrder.NewUploadContractFileLogic(r.Context(), svcCtx)
		resp, err := l.UploadContractFile(r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
