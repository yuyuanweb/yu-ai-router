package errno

type ErrorCode struct {
	Code    int
	Message string
}

var (
	Success        = ErrorCode{Code: 0, Message: "ok"}
	ParamsError    = ErrorCode{Code: 40000, Message: "请求参数错误"}
	NotLoginError  = ErrorCode{Code: 40100, Message: "未登录"}
	NoAuthError    = ErrorCode{Code: 40101, Message: "无权限"}
	TooManyRequest = ErrorCode{Code: 42900, Message: "请求过于频繁"}
	NotFoundError  = ErrorCode{Code: 40400, Message: "请求数据不存在"}
	ForbiddenError = ErrorCode{Code: 40300, Message: "禁止访问"}
	SystemError    = ErrorCode{Code: 50000, Message: "系统内部异常"}
	OperationError = ErrorCode{Code: 50001, Message: "操作失败"}
)
