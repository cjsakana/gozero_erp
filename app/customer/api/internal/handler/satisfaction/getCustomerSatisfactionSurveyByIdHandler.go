package satisfaction

import (
	"net/http"

	"erp/app/customer/api/internal/logic/satisfaction"
	"erp/app/customer/api/internal/svc"
	"erp/app/customer/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetCustomerSatisfactionSurveyByIdHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetCustomerSatisfactionSurveyByIdReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := satisfaction.NewGetCustomerSatisfactionSurveyByIdLogic(r.Context(), svcCtx)
		resp, err := l.GetCustomerSatisfactionSurveyById(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
