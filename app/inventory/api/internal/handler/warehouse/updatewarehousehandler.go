package warehouse

import (
	"net/http"

	"erp/app/inventory/api/internal/logic/warehouse"
	"erp/app/inventory/api/internal/svc"
	"erp/app/inventory/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UpdateWarehouseHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateWarehouseReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := warehouse.NewUpdateWarehouseLogic(r.Context(), svcCtx)
		resp, err := l.UpdateWarehouse(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
