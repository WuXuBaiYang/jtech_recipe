package common

import (
	"errors"
	"server/tool"
)

// SendSMSVerify 发送短信验证码
func SendSMSVerify(phone string) (*string, error) {
	if !tool.VerifyPhoneNumber(phone) {
		return nil, errors.New("手机号校验失败")
	}
	code := phone[len(phone)-4:]
	return &code, nil
}
