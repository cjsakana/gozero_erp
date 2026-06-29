package payroll

import (
	"net/http"

	"erp/app/hr/api/internal/logic/payroll"
	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SearchPayrollsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchPayrollsRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := payroll.NewSearchPayrollsLogic(r.Context(), svcCtx)
		resp, err := l.SearchPayrolls(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
