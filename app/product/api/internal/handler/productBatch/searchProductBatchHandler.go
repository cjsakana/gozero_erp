package productBatch

import (
	"net/http"

	"erp/app/product/api/internal/logic/productBatch"
	"erp/app/product/api/internal/svc"
	"erp/app/product/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SearchProductBatchHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchProductBatchRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := productBatch.NewSearchProductBatchLogic(r.Context(), svcCtx)
		resp, err := l.SearchProductBatch(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
