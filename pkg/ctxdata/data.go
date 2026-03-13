package ctxdata

import "context"

func GetUId(ctx context.Context) string {
	if uid, ok := ctx.Value(Identify).(string); ok {
		return uid
	}
	return ""
}
