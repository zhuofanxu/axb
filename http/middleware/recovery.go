package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/zhuofanxu/axb/http/response"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				var err error
				if e, ok := r.(error); ok {
					err = errors.Wrap(e, "recovered from panic")
				} else {
					err = fmt.Errorf("recovered from panic: %v", r)
				}
				response.NewBaseHandler().Error(c, err)
				c.Abort()
			}
		}()
		c.Next()
	}
}
