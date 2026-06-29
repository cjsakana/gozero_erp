package permission

import (
	"net/http"

	"erp/app/auth/api/internal/logic/permission"
	"erp/app/auth/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetAllPermissionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := permission.NewGetAllPermissionsLogic(r.Context(), svcCtx)
		resp, err := l.GetAllPermissions()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
