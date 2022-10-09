package common

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"reflect"
	"regexp"
	"time"
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
	// 注册小于当前日期判断失败
	if err := v.RegisterValidation("ltToday", verifyLTToday); err != nil {
		panic("校验小于当前时间失败")
	}
	// 注册大于当前日期判断失败
	if err := v.RegisterValidation("gtToday", verifyGTToday); err != nil {
		panic("校验大于当前时间失败")
	}
	// 校验字典类型,传入字典表名，会自动拼接dict头
	if err := v.RegisterValidation("dict", verifyDict); err != nil {
		panic("校验字典失败")
	}
}

// 验证是否为手机号
func verifyPhone(fl validator.FieldLevel) bool {
	regExp := "^1(3\\d|4[5-9]|5[0-35-9]|6[567]|7[0-8]|8\\d|9[0-35-9])\\d{8}$"
	ok, _ := regexp.MatchString(regExp, fl.Field().String())
	return ok
}

// 验证是否为url
func verifyUrl(fl validator.FieldLevel) bool {
	regExp := "^(?:(http|https|ftp):\\/\\/)?((?:[\\w-]+\\.)+[a-z0-9]+)((?:\\/[^/?#]*)+)?(\\?[^#]+)?(#.+)?$"
	ok, _ := regexp.MatchString(regExp, fl.Field().String())
	return ok
}

// 验证是否小于当前日期
func verifyLTToday(fl validator.FieldLevel) bool {
	date := fl.Field().Interface().(time.Time)
	return time.Now().After(date)
}

// 验证是否大于当前日期
func verifyGTToday(fl validator.FieldLevel) bool {
	date := fl.Field().Interface().(time.Time)
	return time.Now().Before(date)
}

// 校验字典类型
func verifyDict(fl validator.FieldLevel) bool {
	db := GetDB()
	var count int64
	var codes []string
	k := fl.Field().Kind()
	if k == reflect.Slice {
		if slice, ok := fl.Field().
			Interface().([]string); ok {
			codes = append(codes, slice...)
		}
	} else if k == reflect.String {
		c := fl.Field().String()
		if len(c) != 0 {
			codes = append(codes, c)
		}
	}
	if len(codes) == 0 {
		return true
	}
	db.Table("sys_dict_"+fl.Param()).
		Where("code in ?", codes).
		Count(&count)
	return int(count) == len(codes)
}
