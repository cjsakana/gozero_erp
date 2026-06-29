package userRole

import (
	"net/http"

	"erp/app/auth/api/internal/logic/userRole"
	"erp/app/auth/api/internal/svc"
	"erp/app/auth/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UpdateUserRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateUserRoleRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := userRole.NewUpdateUserRoleLogic(r.Context(), svcCtx)
		resp, err := l.UpdateUserRole(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
