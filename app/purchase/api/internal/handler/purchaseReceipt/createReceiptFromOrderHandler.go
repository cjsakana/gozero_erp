package purchaseReceipt

import (
	"net/http"

	"erp/app/purchase/api/internal/logic/purchaseReceipt"
	"erp/app/purchase/api/internal/svc"
	"erp/app/purchase/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateReceiptFromOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateReceiptFromOrderReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := purchaseReceipt.NewCreateReceiptFromOrderLogic(r.Context(), svcCtx)
		resp, err := l.CreateReceiptFromOrder(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
