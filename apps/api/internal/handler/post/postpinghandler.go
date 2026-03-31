package post

import (
	"net/http"

	"dicetales.com/apps/api/internal/logic/post"
	"dicetales.com/apps/api/internal/svc"
	"dicetales.com/apps/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Post Module Placeholder
func PostPingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PostPlaceholderReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := post.NewPostPingLogic(r.Context(), svcCtx)
		resp, err := l.PostPing(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
