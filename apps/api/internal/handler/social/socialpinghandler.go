package social

import (
	"net/http"

	"dicetales.com/apps/api/internal/logic/social"
	"dicetales.com/apps/api/internal/svc"
	"dicetales.com/apps/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Social Module Placeholder
func SocialPingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SocialPlaceholderReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := social.NewSocialPingLogic(r.Context(), svcCtx)
		resp, err := l.SocialPing(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
