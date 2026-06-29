package customerCategory

import (
	"net/http"

	"erp/app/customer/api/internal/logic/customerCategory"
	"erp/app/customer/api/internal/svc"
	"erp/app/customer/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UpdateCustomerCategoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateCustomerCategoryReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := customerCategory.NewUpdateCustomerCategoryLogic(r.Context(), svcCtx)
		resp, err := l.UpdateCustomerCategory(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
