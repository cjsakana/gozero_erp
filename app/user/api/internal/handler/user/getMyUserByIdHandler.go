package user

import (
	"net/http"

	"erp/app/user/api/internal/logic/user"
	"erp/app/user/api/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetMyUserByIdHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewGetMyUserByIdLogic(r.Context(), svcCtx)
		resp, err := l.GetMyUserById()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
