package userRole

import (
	"net/http"

	"erp/app/auth/api/internal/logic/userRole"
	"erp/app/auth/api/internal/svc"
	"erp/app/auth/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AddUserRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddUserRoleRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := userRole.NewAddUserRoleLogic(r.Context(), svcCtx)
		resp, err := l.AddUserRole(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
