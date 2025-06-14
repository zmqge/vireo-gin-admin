package utils

import (
	"github.com/mojocn/base64Captcha"
)

var store = base64Captcha.DefaultMemStore

// GenerateCaptcha 生成验证码
func GenerateCaptcha() (string, string, error) {
	// 配置验证码
	driver := base64Captcha.NewDriverDigit(80, 240, 4, 0.7, 80)
	captcha := base64Captcha.NewCaptcha(driver, store)

	// 生成验证码
	id, b64s, _, err := captcha.Generate()
	if err != nil {
		return "", "", err
	}

	return id, b64s, nil
}

// VerifyCaptcha 验证验证码
func VerifyCaptcha(id, answer string) bool {
	return store.Verify(id, answer, true)
}
