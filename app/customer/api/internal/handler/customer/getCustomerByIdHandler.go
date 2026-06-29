package customer

import (
	"net/http"

	"erp/app/customer/api/internal/logic/customer"
	"erp/app/customer/api/internal/svc"
	"erp/app/customer/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetCustomerByIdHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetCustomerByIdReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := customer.NewGetCustomerByIdLogic(r.Context(), svcCtx)
		resp, err := l.GetCustomerById(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
