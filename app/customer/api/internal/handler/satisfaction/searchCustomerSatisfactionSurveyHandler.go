package satisfaction

import (
	"net/http"

	"erp/app/customer/api/internal/logic/satisfaction"
	"erp/app/customer/api/internal/svc"
	"erp/app/customer/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SearchCustomerSatisfactionSurveyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchCustomerSatisfactionSurveyReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := satisfaction.NewSearchCustomerSatisfactionSurveyLogic(r.Context(), svcCtx)
		resp, err := l.SearchCustomerSatisfactionSurvey(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
