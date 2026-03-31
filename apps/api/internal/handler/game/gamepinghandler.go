package game

import (
	"net/http"

	"dicetales.com/apps/api/internal/logic/game"
	"dicetales.com/apps/api/internal/svc"
	"dicetales.com/apps/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Game Module Placeholder
func GamePingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GamePlaceholderReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := game.NewGamePingLogic(r.Context(), svcCtx)
		resp, err := l.GamePing(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
