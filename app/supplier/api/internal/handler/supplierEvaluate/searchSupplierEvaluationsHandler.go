package supplierEvaluate

import (
	"net/http"

	"erp/app/supplier/api/internal/logic/supplierEvaluate"
	"erp/app/supplier/api/internal/svc"
	"erp/app/supplier/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SearchSupplierEvaluationsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchSupplierEvaluationsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := supplierEvaluate.NewSearchSupplierEvaluationsLogic(r.Context(), svcCtx)
		resp, err := l.SearchSupplierEvaluations(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
