package employeeDetail

import (
	"net/http"

	"erp/app/hr/api/internal/logic/employeeDetail"
	"erp/app/hr/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ImportEmployeeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := employeeDetail.NewImportEmployeeLogic(r.Context(), svcCtx)
		resp, err := l.ImportEmployee(r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
