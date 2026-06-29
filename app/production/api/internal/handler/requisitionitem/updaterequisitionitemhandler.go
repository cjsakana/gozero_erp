package requisitionitem

import (
	"net/http"

	"erp/app/production/api/internal/logic/requisitionitem"
	"erp/app/production/api/internal/svc"
	"erp/app/production/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 更新领料单明细
func UpdateRequisitionItemHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateRequisitionItemReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := requisitionitem.NewUpdateRequisitionItemLogic(r.Context(), svcCtx)
		resp, err := l.UpdateRequisitionItem(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
