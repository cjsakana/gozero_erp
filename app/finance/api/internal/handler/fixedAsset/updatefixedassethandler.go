package fixedAsset

import (
	"net/http"

	"erp/app/finance/api/internal/logic/fixedAsset"
	"erp/app/finance/api/internal/svc"
	"erp/app/finance/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UpdateFixedAssetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateFixedAssetReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := fixedAsset.NewUpdateFixedAssetLogic(r.Context(), svcCtx)
		resp, err := l.UpdateFixedAsset(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
