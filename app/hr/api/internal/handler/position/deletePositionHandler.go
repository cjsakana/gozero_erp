package position

import (
	"net/http"

	"erp/app/hr/api/internal/logic/position"
	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DeletePositionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeletePositionRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := position.NewDeletePositionLogic(r.Context(), svcCtx)
		resp, err := l.DeletePosition(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
