package image

import (
	"net/http"

	"erp/app/image/api/internal/logic/image"
	"erp/app/image/api/internal/svc"
	"erp/app/image/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DelImageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DelImageReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := image.NewDelImageLogic(r.Context(), svcCtx)
		resp, err := l.DelImage(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
