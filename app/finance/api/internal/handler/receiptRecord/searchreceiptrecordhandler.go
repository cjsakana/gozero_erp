package receiptRecord

import (
	"net/http"

	"erp/app/finance/api/internal/logic/receiptRecord"
	"erp/app/finance/api/internal/svc"
	"erp/app/finance/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SearchReceiptRecordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchReceiptRecordReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := receiptRecord.NewSearchReceiptRecordLogic(r.Context(), svcCtx)
		resp, err := l.SearchReceiptRecord(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
