package errx

import (
	"errors"
)

// Error 自定义错误结构
type Error struct {
	Code    int    `json:"code"`    // 业务码
	Message string `json:"message"` // 给用户看的提示
	Cause   error  `json:"-"`       // 原始错误 (用于记录日志，不输出到 JSON)
}

// 实现 error 接口
func (e *Error) Error() string {
	return e.Message
}

func (e *Error) WithMsg(msg string) *Error {
	if msg != "" {
		e.Message = msg
	}
	return e
}

// NewError 工厂方法
func NewError(code int, cause error) *Error {
	// 默认使用code对应的错误消息
	return &Error{
		Code:    code,
		Message: GetMessage(code),
		Cause:   cause,
	}
}

// InternalError 创建一个内部服务器错误的 Error
func InternalError(cause error) *Error {
	return NewError(CodeSystemError, cause)
}

// ParamError 创建一个参数错误的 Error
func ParamError(cause error) *Error {
	return NewError(CodeParamError, cause)
}

func UserPasswordError(cause error) *Error {
	return NewError(CodePasswordError, cause)
}

// NotFoundError 创建一个资源未找到的 Error
func NotFoundError(cause error) *Error {
	return NewError(CodeNotFound, cause)
}

// UnauthorizedError 创建一个未授权的 Error
func UnauthorizedError(cause error) *Error {
	return NewError(CodeUnauthorized, cause)
}

// ForbiddenError 创建一个禁止访问的 Error
func ForbiddenError(cause error) *Error {
	return NewError(CodeForbidden, cause)
}

func AlreadyExistsError(cause error) *Error {
	return NewError(CodeAlreadyExists, cause)
}

func CanceledError(cause error) *Error {
	return NewError(CodeCanceled, cause)
}

func TimeoutError(cause error) *Error {
	return NewError(CodeTimout, cause)
}

func IsErrorCode(err error, code int) bool {
	var e *Error
	return errors.As(err, &e) && e.Code == code
}

func UpdateMessageIfCodeMatch(err error, code int, msg string) error {
	var e *Error
	if errors.As(err, &e) && e.Code == code {
		e.Message = msg
		return e
	}
	return err
}
