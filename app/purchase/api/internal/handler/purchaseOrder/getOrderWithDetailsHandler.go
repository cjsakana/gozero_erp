package purchaseOrder

import (
	"net/http"

	"erp/app/purchase/api/internal/logic/purchaseOrder"
	"erp/app/purchase/api/internal/svc"
	"erp/app/purchase/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetOrderWithDetailsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetOrderWithDetailsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := purchaseOrder.NewGetOrderWithDetailsLogic(r.Context(), svcCtx)
		resp, err := l.GetOrderWithDetails(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
