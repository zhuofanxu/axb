package middleware

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/zhuofanxu/axb/errx"
)

func Logger(log *zap.Logger, isProd bool) gin.HandlerFunc {

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		cost := time.Since(start)

		// 准备日志字段
		fields := []zap.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Duration("cost", cost),
		}

		// 根据状态码和错误情况选择日志级别
		statusCode := c.Writer.Status()
		var appErr *errx.Error
		if len(c.Errors) > 0 {
			ginErr := c.Errors.Last()
			err := ginErr.Err

			var validErr *errx.Error
			if errors.As(err, &validErr) {
				appErr = validErr
			} else {
				appErr = errx.InternalError(err)
			}
		}

		if statusCode >= 500 {
			stacktrace := ""
			if appErr != nil {
				fields = append(fields, zap.String("error", appErr.Error()))
				if appErr.Cause != nil {
					// 获取完整的错误信息(包含消息和堆栈)
					lines := strings.Split(fmt.Sprintf("%+v", appErr.Cause), "\n")
					maxBound := 12
					if len(lines) < maxBound {
						maxBound = len(lines)
					}
					stacktrace = strings.Join(lines[:maxBound], "\n") // 取前12行堆栈信息
					if isProd {
						fields = append(fields, zap.String("stacktrace", stacktrace))
					}
				}
			}
			log.Error("request failed", fields...)
			if !isProd && stacktrace != "" {
				fmt.Println(stacktrace)
			}
			return
		}

		if statusCode >= 400 {
			if appErr != nil {
				var err error = appErr
				if appErr.Cause != nil {
					err = appErr.Cause
				}

				fields = append(fields, zap.Error(err))
			}
			log.Warn("request warning", fields...)
			return
		}

		log.Info("request successful", fields...)
	}
}
