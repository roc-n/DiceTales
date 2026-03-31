package meetup

import (
	"net/http"

	"dicetales.com/apps/api/internal/logic/meetup"
	"dicetales.com/apps/api/internal/svc"
	"dicetales.com/apps/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Meetup Module Placeholder
func MeetupPingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MeetupPlaceholderReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := meetup.NewMeetupPingLogic(r.Context(), svcCtx)
		resp, err := l.MeetupPing(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
