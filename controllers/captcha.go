package controllers

import (
	"kuangchi_backend/models"
	"kuangchi_backend/utils"
	"github.com/astaxie/beego/utils/captcha"
	"strings"
	"github.com/astaxie/beego/validation"
)

type CaptchaController struct {
	CommonController
}

// @router /captcha/code [get]
func (u *CaptchaController) Code() {
	code := utils.RandomCaptcha(4)
	u.Ctx.Output.Header("Content-Type", "image/png")
	u.SetSession("code", code)
	intSlice := []byte{}
	for _, char := range code {
		intSlice = append(intSlice, byte(char-'0'))
	}
	img := captcha.NewImage(intSlice, 100, 40)
	img.WriteTo(u.Ctx.ResponseWriter)
}

// @router /captcha/send-captcha [get]
func (u *CaptchaController) SendCaptcha() {
	ident := u.GetString("ident")
	captchaType := u.GetString("type")
	countryCode := u.GetString("country_code")

	if captchaType != "REGISTER" && captchaType != "RESETPASSWORD" {
		u.Error("invalid type")
		return
	}

	if ident == "" {
		u.Error("missing ident")
		return
	}
	valid := validation.Validation{}
	if strings.Contains(ident, "@") {
		if v := valid.Email(ident, "email"); !v.Ok {
			u.Error(100301)
			return
		}

		if captchaType == "REGISTER" && !models.IsUsableEmail(ident) {
			u.Error(100305)
			return
		}

		if captchaType == "RESETPASSWORD" && models.IsUsableEmail(ident) {
			u.Error(100302)
			return
		}

		if err := models.SendEmail(ident, captchaType); err != nil {
			u.Error(err)
			return
		}
	} else {
		if !utils.ValidateMobile(ident) {
			u.Error("Please enter the correct mail address or phone number")
			return
		}

		if captchaType == "REGISTER" && !models.IsUsableMobile(ident) {
			u.Error(100313)
			return
		}

		if captchaType == "RESETPASSWORD" && models.IsUsableMobile(ident) {
			u.Error(100312)
			return
		}
		if captchaType == "RESETPASSWORD" && countryCode == "" {
			user := models.GetUserByMobile(ident)
			if user != nil {
				countryCode = user.CountryCode
			}
		}
		if err := models.SendSMS(ident, captchaType, countryCode); err != nil {
			u.Error(err)
			return
		}
	}
	u.Ok()
}
