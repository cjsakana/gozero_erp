package salesOrder

import (
	"net/http"

	"erp/app/sale/api/internal/logic/salesOrder"
	"erp/app/sale/api/internal/svc"
	"erp/app/sale/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SearchSalesOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchSalesOrderReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := salesOrder.NewSearchSalesOrderLogic(r.Context(), svcCtx)
		resp, err := l.SearchSalesOrder(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
