package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ValidateMobile(field validator.FieldLevel) bool {
	//使用正则表达式验证手机号是否符合要求
	mobile := field.Field().String() //获取手机号
	ok, _ := regexp.MatchString(`^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`, mobile)
	if !ok {
		return false
	}
	return true
}
