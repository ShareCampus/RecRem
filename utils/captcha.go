package utils

import "github.com/mojocn/base64Captcha"

var store = base64Captcha.DefaultMemStore

// CaptchaConfig 验证码结构体
type CaptchaConfig struct {
	Id            string                       `json:"id"`
	CaptchaType   string                       `json:"captcha_type"`
	VerifyValue   string                       `json:"verify_value"`
	DriverAudio   *base64Captcha.DriverAudio   `json:"driver_audio"`
	DriverString  *base64Captcha.DriverString  `json:"driver_string"`
	DriverChinese *base64Captcha.DriverChinese `json:"driver_chinese"`
	DriverMath    *base64Captcha.DriverMath    `json:"driver_math"`
	DriverDigit   *base64Captcha.DriverDigit   `json:"driver_digit"`
}

// CaptchaVerify 校验验证码
func CaptchaVerify(captcha *CaptchaConfig) bool {
	return store.Verify(captcha.Id, captcha.VerifyValue, false)
}
