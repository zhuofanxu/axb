package validation

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/zhuofanxu/axb/errx"
)

func wrapValidationErrors(err error, info interface{}) *errx.Error {
	var valErr validator.ValidationErrors
	if err == io.EOF {
		return errx.NewError(errx.CodeParamError, err).
			WithMsg("Request body is empty. Please provide a valid JSON object.")
	}
	if errors.As(err, &valErr) {
		for _, fieldErr := range valErr {
			// just return the first validation error
			return errx.NewError(errx.CodeParamError, nil).WithMsg(formatValidationError(fieldErr, info))
		}
	}

	return errx.NewError(errx.CodeParamError, err)
}

func formatValidationError(fieldErr validator.FieldError, info interface{}) string {
	var msg string
	switch fieldErr.Tag() {
	case "required":
		msg = "is required"
	case "oneof":
		msg = fmt.Sprintf("must be one of [%s]", fieldErr.Param())
	case "min":
		msg = fmt.Sprintf("must be at least %s", fieldErr.Param())
	case "max":
		msg = fmt.Sprintf("must be at most %s", fieldErr.Param())
	case "len":
		msg = fmt.Sprintf("must be exactly %s characters long", fieldErr.Param())
	case "gte":
		msg = fmt.Sprintf("must be greater than or equal to %s", fieldErr.Param())
	case "gt":
		msg = fmt.Sprintf("must be greater than %s", fieldErr.Param())
	case "lte":
		msg = fmt.Sprintf("must be less than or equal to %s", fieldErr.Param())
	case "lt":
		msg = fmt.Sprintf("must be less than %s", fieldErr.Param())
	case "email":
		msg = "must be a valid email address"
	default:
		if fieldErr.Param() != "" {
			msg = fmt.Sprintf("is invalid (rule: %s, expected: %s)", fieldErr.Tag(), fieldErr.Param())
		} else {
			msg = fmt.Sprintf("is invalid (rule: %s)", fieldErr.Tag())
		}
	}

	base := fmt.Sprintf("Field '%s' %s.", fieldErr.Field(), msg)
	if info == nil {
		return base
	}
	infoText := strings.TrimSpace(fmt.Sprint(info))
	if infoText == "" {
		return base
	}
	return fmt.Sprintf("%s: %s", infoText, base)
}

func BindUrlParam(param interface{}, c *gin.Context) error {
	if err := c.ShouldBindUri(param); err != nil {
		return wrapValidationErrors(err, "Invalid path parameters")
	}
	return nil
}

func BindQueryParam(query interface{}, c *gin.Context) error {
	if err := c.ShouldBindQuery(query); err != nil {
		return wrapValidationErrors(err, "Invalid query parameters")
	}
	return nil
}

func BindJsonBody(body interface{}, c *gin.Context) error {
	if err := c.ShouldBindJSON(body); err != nil {
		return wrapValidationErrors(err, "Invalid JSON request body")
	}
	return nil
}
