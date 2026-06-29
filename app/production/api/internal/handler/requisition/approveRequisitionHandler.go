package requisition

import (
	"net/http"

	"erp/app/production/api/internal/logic/requisition"
	"erp/app/production/api/internal/svc"
	"erp/app/production/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 审批领料单
func ApproveRequisitionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ApproveRequisitionReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := requisition.NewApproveRequisitionLogic(r.Context(), svcCtx)
		resp, err := l.ApproveRequisition(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
