package user

import (
	"net/http"

	"dicetales.com/apps/api/internal/logic/user"
	"dicetales.com/apps/api/internal/svc"
	"dicetales.com/apps/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 刷新Token
func RefreshTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RefreshTokenReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewRefreshTokenLogic(r.Context(), svcCtx)
		resp, err := l.RefreshToken(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
