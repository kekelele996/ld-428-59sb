package util

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var cnMobilePattern = regexp.MustCompile(`^1[3-9]\d{9}$`)

// RegisterCustomValidations 注册业务自定义校验标签（如 cn_mobile 国内手机号）。
func RegisterCustomValidations() error {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return errors.New("binding validator engine unavailable")
	}
	return v.RegisterValidation("cn_mobile", func(fl validator.FieldLevel) bool {
		return cnMobilePattern.MatchString(fl.Field().String())
	})
}

// TranslateError 将参数校验错误翻译为可读的中文提示。
func TranslateError(err error) string {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		first := ve[0]
		return fmt.Sprintf("字段 %s 校验失败：%s", first.Field(), humanizeTag(first.Tag()))
	}
	return err.Error()
}

func humanizeTag(tag string) string {
	switch tag {
	case "required":
		return "不能为空"
	case "min":
		return "长度或值过小"
	case "max":
		return "长度或值超出限制"
	case "len":
		return "长度不正确"
	case "oneof":
		return "取值不在允许范围"
	case "email":
		return "格式不正确"
	default:
		return tag
	}
}
