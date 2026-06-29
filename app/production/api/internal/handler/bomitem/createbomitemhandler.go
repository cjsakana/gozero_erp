package bomitem

import (
	"net/http"

	"erp/app/production/api/internal/logic/bomitem"
	"erp/app/production/api/internal/svc"
	"erp/app/production/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 创建BOM明细
func CreateBomItemHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateBomItemReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := bomitem.NewCreateBomItemLogic(r.Context(), svcCtx)
		resp, err := l.CreateBomItem(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
