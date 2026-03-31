package ctxdata

import "context"

const Identify string = "dicetales.uid"

func GetUId(ctx context.Context) string {
	if uid, ok := ctx.Value(Identify).(string); ok {
		return uid
	}
	return ""
}
