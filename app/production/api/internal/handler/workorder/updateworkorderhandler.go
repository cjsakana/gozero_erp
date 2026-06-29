package workorder

import (
	"net/http"

	"erp/app/production/api/internal/logic/workorder"
	"erp/app/production/api/internal/svc"
	"erp/app/production/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 更新生产工单
func UpdateWorkOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateWorkOrderReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := workorder.NewUpdateWorkOrderLogic(r.Context(), svcCtx)
		resp, err := l.UpdateWorkOrder(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
