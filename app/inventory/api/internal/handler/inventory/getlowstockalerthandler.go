package inventory

import (
	"net/http"

	"erp/app/inventory/api/internal/logic/inventory"
	"erp/app/inventory/api/internal/svc"
	"erp/app/inventory/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetLowStockAlertHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetLowStockAlertReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := inventory.NewGetLowStockAlertLogic(r.Context(), svcCtx)
		resp, err := l.GetLowStockAlert(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
