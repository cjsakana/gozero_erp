package userRole

import (
	"net/http"

	"erp/app/auth/api/internal/logic/userRole"
	"erp/app/auth/api/internal/svc"
	"erp/app/auth/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DelUserRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DelUserRoleRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := userRole.NewDelUserRoleLogic(r.Context(), svcCtx)
		resp, err := l.DelUserRole(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
