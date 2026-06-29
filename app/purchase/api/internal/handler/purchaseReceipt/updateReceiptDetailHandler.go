package purchaseReceipt

import (
	"net/http"

	"erp/app/purchase/api/internal/logic/purchaseReceipt"
	"erp/app/purchase/api/internal/svc"
	"erp/app/purchase/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UpdateReceiptDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateReceiptDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := purchaseReceipt.NewUpdateReceiptDetailLogic(r.Context(), svcCtx)
		resp, err := l.UpdateReceiptDetail(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
