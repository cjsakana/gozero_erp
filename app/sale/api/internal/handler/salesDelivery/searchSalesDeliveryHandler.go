package salesDelivery

import (
	"net/http"

	"erp/app/sale/api/internal/logic/salesDelivery"
	"erp/app/sale/api/internal/svc"
	"erp/app/sale/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SearchSalesDeliveryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchSalesDeliveryReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := salesDelivery.NewSearchSalesDeliveryLogic(r.Context(), svcCtx)
		resp, err := l.SearchSalesDelivery(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
