package purchaseRequisition

import (
	"net/http"

	"erp/app/purchase/api/internal/logic/purchaseRequisition"
	"erp/app/purchase/api/internal/svc"
	"erp/app/purchase/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetRequisitionWithDetailsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetRequisitionWithDetailsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := purchaseRequisition.NewGetRequisitionWithDetailsLogic(r.Context(), svcCtx)
		resp, err := l.GetRequisitionWithDetails(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
