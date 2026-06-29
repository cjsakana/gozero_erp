package requisition

import (
	"net/http"

	"erp/app/production/api/internal/logic/requisition"
	"erp/app/production/api/internal/svc"
	"erp/app/production/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 更新领料单
func UpdateRequisitionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateRequisitionReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := requisition.NewUpdateRequisitionLogic(r.Context(), svcCtx)
		resp, err := l.UpdateRequisition(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
