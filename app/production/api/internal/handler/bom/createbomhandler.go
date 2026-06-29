package bom

import (
	"net/http"

	"erp/app/production/api/internal/logic/bom"
	"erp/app/production/api/internal/svc"
	"erp/app/production/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 创建BOM
func CreateBomHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateBomReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := bom.NewCreateBomLogic(r.Context(), svcCtx)
		resp, err := l.CreateBom(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
