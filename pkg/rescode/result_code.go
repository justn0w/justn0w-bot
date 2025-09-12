package rescode

type ErrorCode struct {
	Code    int
	Message string
}

var (
	ErrUserNotFound = ErrorCode{10001, "用户不存在"}

	// TokenInvalid Token相关错误 20000-29999
	TokenInvalid    = ErrorCode{20001, "无效的token"}
	TokenExpired    = ErrorCode{20002, "token已过期"}
	TokenAuthEmpty  = ErrorCode{20003, "token不能为空"}
	TokenAuthFormat = ErrorCode{20004, "token格式错误"}
)
