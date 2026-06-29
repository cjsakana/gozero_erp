package position

import (
	"net/http"

	"erp/app/hr/api/internal/logic/position"
	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetPositionEmployeesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetPositionEmployeesRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := position.NewGetPositionEmployeesLogic(r.Context(), svcCtx)
		resp, err := l.GetPositionEmployees(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
