package productBatch

import (
	"net/http"

	"erp/app/product/api/internal/logic/productBatch"
	"erp/app/product/api/internal/svc"
	"erp/app/product/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AddProductBatchHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddProductBatchRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := productBatch.NewAddProductBatchLogic(r.Context(), svcCtx)
		resp, err := l.AddProductBatch(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
