package salesOrder

import (
	"fmt"
	"net/http"

	"erp/app/sale/api/internal/logic/salesOrder"
	"erp/app/sale/api/internal/svc"
	"erp/app/sale/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AddSalesOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddSalesOrderReq
		if err := httpx.Parse(r, &req); err != nil {
			fmt.Println("11111111", req)
			fmt.Println("222222222", err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := salesOrder.NewAddSalesOrderLogic(r.Context(), svcCtx)
		resp, err := l.AddSalesOrder(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
