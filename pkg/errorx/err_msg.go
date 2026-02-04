package errorx

var codeText = map[int]string{
	SERVER_COMMON_ERROR: "服务器异常",
	REQUEST_PARAM_ERROR: "请求参数有误",
	DB_ERROR:            "数据库请求出错",
}

func ErrMsg(code int) string {
	if msg, ok := codeText[code]; ok {
		return msg
	}

	// 默认返回服务器异常
	return codeText[SERVER_COMMON_ERROR]
}
