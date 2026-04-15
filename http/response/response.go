package response

import (
	"errors"
	"net/http"

	"github.com/zhuofanxu/axb/errx"

	"github.com/gin-gonic/gin"
)

// Response response is the standard 'response' structure for API responses.
type Response[T any] struct {
	Code    int    `json:"code"`    // 业务码
	Message string `json:"message"` // 响应消息
	Data    T      `json:"result"`  // 响应数据
	Status  int    `json:"status"`  // HTTP 状态码
}

type BaseHandler struct{}

func NewBaseHandler() *BaseHandler {
	return &BaseHandler{}
}

// Success 成功响应
func (h *BaseHandler) Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &Response[any]{
		Code:    errx.CodeSuccess,
		Message: errx.GetMessage(errx.CodeSuccess),
		Data:    data,
		Status:  http.StatusOK,
	})
}

// Error 错误响应
func (h *BaseHandler) Error(c *gin.Context, err error) {
	// default error information
	code := errx.CodeSystemError
	msg := errx.GetMessage(errx.CodeSystemError) // 不直接暴露系统错误详情
	httpStatus := http.StatusInternalServerError
	var data interface{}

	// check and extract wrapped error information
	var wrappedErr *errx.Error
	if errors.As(err, &wrappedErr) {
		code = wrappedErr.Code
		msg = wrappedErr.Message
		httpStatus = getHttpStatusByCode(code)
	}

	_ = c.Error(err)
	c.JSON(httpStatus, &Response[any]{
		Code:    code,
		Message: msg,
		Data:    data,
		Status:  httpStatus,
	})
}

func getHttpStatusByCode(code int) int {
	switch code {
	case errx.CodeSystemError:
		return http.StatusInternalServerError
	case errx.CodeParamError:
		return http.StatusBadRequest
	case errx.CodeUnauthorized:
		return http.StatusUnauthorized
	case errx.CodeForbidden:
		return http.StatusForbidden
	case errx.CodeNotFound:
		return http.StatusNotFound
	case errx.CodeAlreadyExists:
		return http.StatusConflict
	case errx.CodeTimout, errx.CodeCanceled:
		return http.StatusRequestTimeout
	case errx.CodeInvalidToken, errx.CodeTokenExpired:
		return http.StatusUnauthorized
	case errx.CodeNoChanged:
		// 无数据变更也返回200
		return http.StatusOK
	default:
		// 业务级错误返回400
		return http.StatusBadRequest
	}
}
