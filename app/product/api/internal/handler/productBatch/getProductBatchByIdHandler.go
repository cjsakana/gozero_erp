package productBatch

import (
	"net/http"

	"erp/app/product/api/internal/logic/productBatch"
	"erp/app/product/api/internal/svc"
	"erp/app/product/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetProductBatchByIdHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetProductBatchByIdRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := productBatch.NewGetProductBatchByIdLogic(r.Context(), svcCtx)
		resp, err := l.GetProductBatchById(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
