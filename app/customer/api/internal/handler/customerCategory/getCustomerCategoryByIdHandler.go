package customerCategory

import (
	"net/http"

	"erp/app/customer/api/internal/logic/customerCategory"
	"erp/app/customer/api/internal/svc"
	"erp/app/customer/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetCustomerCategoryByIdHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetCustomerCategoryByIdReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := customerCategory.NewGetCustomerCategoryByIdLogic(r.Context(), svcCtx)
		resp, err := l.GetCustomerCategoryById(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
