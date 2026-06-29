package fixedAsset

import (
	"net/http"

	"erp/app/finance/api/internal/logic/fixedAsset"
	"erp/app/finance/api/internal/svc"
	"erp/app/finance/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetFixedAssetByIdHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetFixedAssetByIdReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := fixedAsset.NewGetFixedAssetByIdLogic(r.Context(), svcCtx)
		resp, err := l.GetFixedAssetById(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
