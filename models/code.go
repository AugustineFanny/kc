package models

import (
	"kuangchi_backend/result"
	"kuangchi_backend/utils"
)

func SendEmail(email, captchaType string) error {
	if _, err := utils.RedisClient.GetString("HOLD:" + captchaType + ":" + email); err == nil {
		return result.ErrCode(100112)
	}

	captcha := utils.RandomCaptcha(5)

	if !utils.CaptchaSet(email, captcha, captchaType) {
		return result.ErrCode(100304)
	}

	utils.RedisClient.Set("HOLD:" + captchaType + ":" + email, "hold", 50)

	go utils.SendEmail(email, captcha)
	return nil
}

func SendSMS(mobile, captchaType, countryCode string) error {
	if _, err := utils.RedisClient.GetString("HOLD:" + captchaType + ":" + countryCode + mobile); err == nil {
		return result.ErrCode(100112)
	}

	captcha := utils.RandomCaptcha(5)

	if !utils.CaptchaSet(countryCode + mobile, captcha, captchaType) {
		return result.ErrCode(100304)
	}

	utils.RedisClient.Set("HOLD:" + captchaType + ":" + countryCode + mobile, "hold", 50)

	go utils.SendSMS(countryCode + mobile, captcha)
	return nil
}
