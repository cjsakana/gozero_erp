package attendanceReplenish

import (
	"net/http"

	"erp/app/hr/api/internal/logic/attendanceReplenish"
	"erp/app/hr/api/internal/svc"
	"erp/app/hr/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SearchAttendanceReplenishHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SearchAttendanceReplenishRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := attendanceReplenish.NewSearchAttendanceReplenishLogic(r.Context(), svcCtx)
		resp, err := l.SearchAttendanceReplenish(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
