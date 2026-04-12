package ctxdata

import "context"

type CtxKey string

const Identify CtxKey = "dicetales.uid"

func GetUId(ctx context.Context) string {
	if uid, ok := ctx.Value(Identify).(string); ok {
		return uid
	}
	return ""
}
