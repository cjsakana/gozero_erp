package department

import (
	"net/http"

	"erp/app/hr/api/internal/logic/department"
	"erp/app/hr/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ImportDepartmentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := department.NewImportDepartmentLogic(r.Context(), svcCtx)
		resp, err := l.ImportDepartment(r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
