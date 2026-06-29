package employeeDetail

import (
	"net/http"

	"erp/app/hr/api/internal/logic/employeeDetail"
	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AdjustSalaryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdjustSalaryRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := employeeDetail.NewAdjustSalaryLogic(r.Context(), svcCtx)
		resp, err := l.AdjustSalary(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
