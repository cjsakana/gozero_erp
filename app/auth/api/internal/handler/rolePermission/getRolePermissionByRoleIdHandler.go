package rolePermission

import (
	"fmt"
	"net/http"

	"erp/app/auth/api/internal/logic/rolePermission"
	"erp/app/auth/api/internal/svc"
	"erp/app/auth/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetRolePermissionByRoleIdHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetRolePermissionByRoleIdRequest
		if err := httpx.Parse(r, &req); err != nil {
			fmt.Println("123", err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := rolePermission.NewGetRolePermissionByRoleIdLogic(r.Context(), svcCtx)
		resp, err := l.GetRolePermissionByRoleId(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
