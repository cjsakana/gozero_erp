package salesDelivery

import (
	"net/http"

	"erp/app/sale/api/internal/logic/salesDelivery"
	"erp/app/sale/api/internal/svc"
	"erp/app/sale/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetSalesDeliveryByIdHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetSalesDeliveryByIdReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := salesDelivery.NewGetSalesDeliveryByIdLogic(r.Context(), svcCtx)
		resp, err := l.GetSalesDeliveryById(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
