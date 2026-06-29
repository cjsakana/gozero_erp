package salaryPayment

import (
	"net/http"

	"erp/app/finance/api/internal/logic/salaryPayment"
	"erp/app/finance/api/internal/svc"
	"erp/app/finance/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SearchSalaryPaymentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchSalaryPaymentReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := salaryPayment.NewSearchSalaryPaymentLogic(r.Context(), svcCtx)
		resp, err := l.SearchSalaryPayment(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
