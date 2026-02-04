package errorx

const (
	SERVER_COMMON_ERROR = 100001 // 服务器内部错误
	REQUEST_PARAM_ERROR = 100002 // 请求参数错误
	DB_ERROR            = 100003 // 数据库错误

	RESOURCE_ALREADY_EXISTS = 200001 // 资源已存在
	BUSINESS_LOGIC_ERROR    = 200002 // 业务逻辑错误
)
