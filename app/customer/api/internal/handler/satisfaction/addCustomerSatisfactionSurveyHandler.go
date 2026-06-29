package satisfaction

import (
	"net/http"

	"erp/app/customer/api/internal/logic/satisfaction"
	"erp/app/customer/api/internal/svc"
	"erp/app/customer/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AddCustomerSatisfactionSurveyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddCustomerSatisfactionSurveyReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := satisfaction.NewAddCustomerSatisfactionSurveyLogic(r.Context(), svcCtx)
		resp, err := l.AddCustomerSatisfactionSurvey(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
