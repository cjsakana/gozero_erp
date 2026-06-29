package salesDelivery

import (
	"net/http"

	"erp/app/sale/api/internal/logic/salesDelivery"
	"erp/app/sale/api/internal/svc"
	"erp/app/sale/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func OutboundHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OutboundReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := salesDelivery.NewOutboundLogic(r.Context(), svcCtx)
		resp, err := l.Outbound(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
