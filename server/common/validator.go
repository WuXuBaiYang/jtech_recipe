package common

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"regexp"
)

// InitValidator 初始化自定义验证方法
func InitValidator() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		panic("验证引擎初始化失败")
	}
	// 注册手机号验证方法
	if err := v.RegisterValidation("phone", verifyPhone); err != nil {
		panic("手机号校验失败")
	}
	// 注册url验证方法
	if err := v.RegisterValidation("url", verifyUrl); err != nil {
		panic("url校验失败")
	}
}

// 验证是否为手机号
func verifyPhone(fl validator.FieldLevel) bool {
	regExp := "^1(3\\d|4[5-9]|5[0-35-9]|6[567]|7[0-8]|8\\d|9[0-35-9])\\d{8}$"
	ok, _ := regexp.MatchString(regExp, fl.Param())
	return ok
}

// 验证是否为url
func verifyUrl(fl validator.FieldLevel) bool {
	regExp := "^(?=^.{3,255}$)[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+$"
	ok, _ := regexp.MatchString(regExp, fl.Param())
	return ok
}
