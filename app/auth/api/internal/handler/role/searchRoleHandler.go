package role

import (
	"net/http"

	"erp/app/auth/api/internal/logic/role"
	"erp/app/auth/api/internal/svc"
	"erp/app/auth/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SearchRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchRoleRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := role.NewSearchRoleLogic(r.Context(), svcCtx)
		resp, err := l.SearchRole(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
