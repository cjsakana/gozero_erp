package department

import (
	"net/http"

	"erp/app/hr/api/internal/logic/department"
	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetDepartmentTreeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetDepartmentTreeRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := department.NewGetDepartmentTreeLogic(r.Context(), svcCtx)
		resp, err := l.GetDepartmentTree(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
