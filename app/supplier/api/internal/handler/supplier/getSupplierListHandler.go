package supplier

import (
	"net/http"

	"erp/app/supplier/api/internal/logic/supplier"
	"erp/app/supplier/api/internal/svc"
	"erp/app/supplier/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetSupplierListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetSupplierListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := supplier.NewGetSupplierListLogic(r.Context(), svcCtx)
		resp, err := l.GetSupplierList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
