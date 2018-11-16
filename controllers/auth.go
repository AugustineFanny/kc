package controllers

import (
	"encoding/json"
	"kuangchi_backend/models"
	"kuangchi_backend/result"
	"kuangchi_backend/utils"
	"strings"
)

type AuthController struct {
	CommonController
}

// @router /register-captcha [post]
func (u *AuthController) RegisterCaptcha() {
	var userDto models.UserDto
	json.Unmarshal(u.Ctx.Input.RequestBody, &userDto)
	captcha := utils.CaptchaGet(userDto.CountryCode + userDto.Ident, "REGISTER")
	if captcha == "" {
		u.Error(100307)
		return
	}
	if captcha != userDto.Captcha {
		u.Error(100308)
		return
	}
	u.Ok()
}

// @router /register [post]
func (u *AuthController) Register() {
	var userDto models.UserDto
	json.Unmarshal(u.Ctx.Input.RequestBody, &userDto)

	captcha := utils.CaptchaGet(userDto.CountryCode + userDto.Ident, "REGISTER")
	if captcha == "" {
		u.Error(100307)
		return
	}
	if captcha != userDto.Captcha {
		u.Error(100308)
		return
	}

	if !utils.PasswordStrength(userDto.Password) {
		u.Error("The password must contain letters and numbers")
		return
	}

	if userDto.Password != userDto.CheckPassword {
		u.Error(100318)
		return
	}

	user := models.User{Username: userDto.Username, Password: userDto.Password}
	var err error
	if strings.Contains(userDto.Ident, "@") {
		user.Email = userDto.Ident
		err = models.CreateUserByEmail(&user, userDto.InviteCode)
	} else {
		user.Mobile = userDto.Ident
		user.CountryCode = userDto.CountryCode
		err = models.CreateUserByMobile(&user, userDto.InviteCode)
	}

	if err != nil {
		u.Error(err)
		return
	}
	u.Ok()
}

// @router /activate [get]
func (u *AuthController) Activate() {
	key := u.GetString("key")
	err := models.Activate(key)
	if err != nil {
		u.Error(err)
		return
	}
	u.Ctx.Redirect(302, "/")
}

// @router /reset-password [post]
func (u *AuthController) ResetPassword() {
	var userDto models.UserDto
	json.Unmarshal(u.Ctx.Input.RequestBody, &userDto)

	countryCode := ""
	if !strings.Contains(userDto.Ident, "@") {
		//手机号时需要取countryCode
		user := models.GetUserByMobile(userDto.Ident)
		if user != nil {
			countryCode = user.CountryCode
		}
	}
	captcha := utils.CaptchaGet(countryCode + userDto.Ident, "RESETPASSWORD")
	if captcha == "" {
		u.Error(100307)
		return
	}
	if captcha != userDto.Captcha {
		u.Error(100308)
		return
	}

	if !utils.PasswordStrength(userDto.NewPassword) {
		u.Error("The password must contain letters and numbers")
		return
	}

	if userDto.NewPassword != userDto.CheckPassword {
		u.Error(100318)
		return
	}

	var err error
	if strings.Contains(userDto.Ident, "@") {
		if models.IsUsableEmail(userDto.Ident) {
			u.Error(100302)
			return
		}
		err = models.ResetPasswordByEmail(userDto.Ident, userDto.NewPassword)
	} else {
		if models.IsUsableMobile(userDto.Ident) {
			u.Error(100312)
			return
		}
		err = models.ResetPasswordByMobile(userDto.Ident, userDto.NewPassword)
	}

	if err != nil {
		u.Error(err)
		return
	}
	utils.CaptchaDel(userDto.Ident, "RESETPASSWORD")
	u.Ok()
}

// @router /login [post]
func (u *AuthController) Login() {
	var userDto models.UserDto
	json.Unmarshal(u.Ctx.Input.RequestBody, &userDto)

	if userDto.Code != u.GetSession("code") {
		u.Error(100101)
		return
	}

	if userDto.Ident == "" {
		u.Error("missing ident")
		return
	}

	var user *models.User
	var err error
	if strings.Contains(userDto.Ident, "@") {
		user, err = models.UserSignInByEmail(userDto.Ident, userDto.Password)
	} else if userDto.Ident[0] < '0' || userDto.Ident[0] > '9' {
		user, err = models.UserSignInByUsername(userDto.Ident, userDto.Password)
	} else {
		user, err = models.UserSignInByMobile(userDto.Ident, userDto.Password)
	}

	if err != nil {
		u.Error(err)
		return
	}
	u.DelSession("code")
	token := utils.CreateToken(user.Username)
	u.Ctx.ResponseWriter.Header().Add("Authorization", "Bearer " + token)
	u.Data["json"] = result.Success(user)
	u.ServeJSON()
}

// @router /logout [get]
func (u *AuthController) Logout() {
	u.DelSession("user")
	u.Ok()
}
