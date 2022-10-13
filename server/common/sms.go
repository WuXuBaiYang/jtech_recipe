package common

// SendSMSVerify 发送短信验证码
func SendSMSVerify(phone string, code string) error {
	// 如果是调试模式，则不发送验证码
	return nil
}
