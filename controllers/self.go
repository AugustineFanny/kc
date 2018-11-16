package controllers

import (
	"encoding/json"
	"kuangchi_backend/models"
	"kuangchi_backend/result"
	"kuangchi_backend/utils"
	"strings"
)

type SelfController struct {
	CommonController
}

// @router /self/safe [get]
func (u *SelfController) Safe() {
	type data struct {
		models.User
	}
	user := u.GetUser()
	d := data{User: *user}
	u.Data["json"] = result.Success(d)
	u.ServeJSON()
}

// @router /self/info [get]
func (u *SelfController) Info() {
	user := u.GetUser()
	u.Data["json"] = result.Success(models.GetInfo(user))
	u.ServeJSON()
}

// @router /self/avatar [post]
func (u *SelfController) SetAvatar() {
	user := u.GetUser()
	avatar, err := u.SaveImg(user.Id, utils.AvatarImgPath, "avatar", true)
	if err != nil {
		u.Error(err)
		return
	}
	user.Avatar = avatar
	if err := models.SetAvatar(user); err != nil {
		u.Error(err)
		return
	}
	u.Ok()
}

// @router /self/change-password [post]
func (u *SelfController) ChangePassword() {
	user := u.GetUser()

	var dto models.UserDto
	json.Unmarshal(u.Ctx.Input.RequestBody, &dto)

	if !utils.PasswordStrength(dto.NewPassword) {
		u.Error("The password must contain letters and numbers")
		return
	}

	if dto.NewPassword != dto.CheckPassword {
		u.Error(100318)
		return
	}

	err := models.ChangeUserPassword(user, dto.Password, dto.NewPassword)
	if err != nil {
		u.Error(err)
		return
	}
	message := "修改登录密码成功"
	messageEn := "Successfully changed the password."
	messageKo := "로그인 비밀번호를 성공적으로 변경하었습니다."
	messageJp := "パスワード変更成功"
	models.SetMessage(user.Id, message, messageEn,  messageKo, messageJp, "")
	u.Ok()
}

// @router /self/name-auth [get]
func (u *SelfController) RealName() {
	user := u.GetUser()
	res := map[string]interface{}{
		"role": user.Role,
	}
	if user.Verified("role") {
		realName := models.GetRealName(user)
		if realName != nil {
			res["country"] = realName.Country
			res["credential_type"] = realName.CredentialType
			res["name"] = realName.Name
			res["start_date"] = realName.StartDate
			res["end_date"] = realName.EndDate
			res["card"] = utils.ReplaceAsterisk(realName.Card)
		}
	}
	u.Data["json"] = result.Success(res)
	u.ServeJSON()
}

// @router /self/name-auth [post]
func (u *SelfController) NameAuth() {
	user := u.GetUser()
	if models.ExistRealName(user.Id) {
		u.Error(100104)
		return
	}
	realName := models.RealName{Status: 1}
	realName.Country = u.GetString("country")
	realName.CredentialType = u.GetString("credential_type")
	realName.Name = u.GetString("name")
	card := u.GetString("card")
	if realName.CredentialType == "IdCard" && !utils.CheckCard(card) {
		u.Error(100309)
		return
	}
	realName.Card = card

	realName.StartDate = u.GetString("start_date")
	realName.EndDate = u.GetString("end_date")
	var err error
	realName.CardFront, err = u.SaveImg(user.Id, utils.CardImgPath, "card_front", true)
	if err != nil {
		u.Error(err)
		return
	}
	realName.CardBack, err = u.SaveImg(user.Id, utils.CardImgPath, "card_back", true)
	if err != nil {
		u.Error(err)
		return
	}
	realName.CardHold, err = u.SaveImg(user.Id, utils.CardImgPath, "card_hold", true)
	if err != nil {
		u.Error(err)
		return
	}

	if err := models.NameAuth(user, realName); err != nil {
		u.Error(err)
		return
	}
	message := "提交身份验证成功,等待审核中"
	messageEn := "Successfully submitted ID authentication, waiting for system approval."
	messageKo := "신분 인증에 성공하셨습니다. 심사중입니다."
	messageJp := "個人情報提出成功、現在確認中"
	models.SetMessage(user.Id, message, messageEn,  messageKo, messageJp, "")
	u.Ok()
}

// @router /self/kyc [get]
func (u *SelfController) GetKyc() {
	user := u.GetUser()
	res := map[string]interface{}{
		"kyc": user.Kyc,
	}
	if user.Verified("kyc") {
		kyc := models.GetKyc(user)
		if kyc != nil {
			res["name"] = kyc.Name
			res["email"] = kyc.Email
			res["mobile"] = utils.ReplaceAsterisk(kyc.Mobile)
			res["birthday"] = kyc.Birthday
			res["country"] = kyc.Country
			res["province"] = kyc.Province
			res["city"] = kyc.City
			res["street"] = kyc.Street
			res["post_code"] = kyc.PostCode
			res["identity_doc"] = kyc.IdentityDoc
		}
	}
	u.Data["json"] = result.Success(res)
	u.ServeJSON()
}

// @router /self/kyc [post]
func (u *SelfController) Kyc() {
	user := u.GetUser()
	if models.ExistKyc(user.Id) {
		u.Error(100104)
		return
	}
	form := models.Kyc{Status: 1}
	form.Name = u.GetString("name")
	form.Email = u.GetString("email")
	form.Username = u.GetString("username")
	form.Mobile = u.GetString("mobile")
	form.Birthday = u.GetString("birthday")
	form.Country = u.GetString("country")
	form.Province = u.GetString("province")
	form.City = u.GetString("city")
	form.Street = u.GetString("street")
	form.PostCode = u.GetString("post_code")
	form.IdentityDoc = u.GetString("identity_document")
	form.FundsSource = u.GetString("funds_source")

	var err error
	form.IdFront, err = u.SaveImg(user.Id, utils.KycImgPath, "photo_id_front", true)
	if err != nil {
		u.Error(err)
		return
	}
	form.IdBack, err = u.SaveImg(user.Id, utils.KycImgPath, "photo_id_back", true)
	if err != nil {
		u.Error(err)
		return
	}
	form.IdHold, err = u.SaveImg(user.Id, utils.KycImgPath, "photo_id_hold", true)
	if err != nil {
		u.Error(err)
		return
	}
	form.BankAccount, err = u.SaveImg(user.Id, utils.KycImgPath, "bank_account_photo", true)
	if err != nil {
		u.Error(err)
		return
	}
	form.Passport, err = u.SaveImg(user.Id, utils.KycImgPath, "passport", false)
	if err != nil {
		u.Error(err)
		return
	}
	form.Utility, err = u.SaveImg(user.Id, utils.KycImgPath, "utility", false)
	if err != nil {
		u.Error(err)
		return
	}
	form.Driver, err = u.SaveImg(user.Id, utils.KycImgPath, "driver_license", false)
	if err != nil {
		u.Error(err)
		return
	}
	form.TaxBill, err = u.SaveImg(user.Id, utils.KycImgPath, "tax_bill", false)
	if err != nil {
		u.Error(err)
		return
	}
	form.BankStatements, err = u.SaveImg(user.Id, utils.KycImgPath, "bank_statements", false)
	if err != nil {
		u.Error(err)
		return
	}
	form.Other, err = u.SaveImg(user.Id, utils.KycImgPath, "other", false)
	if err != nil {
		u.Error(err)
		return
	}
	if err := models.SetKYC(user, form); err != nil {
		u.Error(err)
		return
	}
	message := "提交KYC信息成功,等待审核中"
	messageEn := "Successfully submitted KYC, waiting for system approval."
	messageKo := "KYC 인증에 성공하셨습니다. 심사중입니다."
	messageJp := "KYC提出成功、現在確認中"
	models.SetMessage(user.Id, message, messageEn,  messageKo, messageJp, "")
	u.Ok()
}

// @router /self/two-factor [get]
func (u *SelfController) TfUrl() {
	user := u.GetUser()
	if user.TfSecret != "" {
		u.Error(100104)
		return
	}
	secretKey := utils.RandomString(20)
	utils.CaptchaSet(user.Username, secretKey, "TWOFACTOR")
	key, url := utils.DefaultOptions.Url("kuangchi_backend-" + user.Username, secretKey)
	resp := map[string]string{
		"key": key,
		"url": url,
	}
	u.Data["json"] = result.Success(resp)
	u.ServeJSON()
}

// @router /self/two-factor [post]
func (u *SelfController) TfInit() {
	var form models.TfCodeDto
	json.Unmarshal(u.Ctx.Input.RequestBody, &form)

	user := u.GetUser()
	if user.TfSecret != "" {
		u.Error(100104)
		return
	}
	secretKey := utils.CaptchaGet(user.Username, "TWOFACTOR")
	if utils.DefaultOptions.Authenticate(secretKey, form.TfCode) {
		if models.SetTwoFactor(user, secretKey, true) {
			models.SetMessage(user.Id, "设置双重认证成功", "", "", "", "")
			u.Ok()
		} else {
			u.Error(100102)
		}
	} else {
		u.Error(100106)
	}
}

// @router /self/two-factor/d [post]
func (u *SelfController) TfDelete() {
	var form models.TfCodeDto
	json.Unmarshal(u.Ctx.Input.RequestBody, &form)

	user := u.GetUser()
	if user.TfOpened == false {
		u.Ok()
		return
	}
	if utils.DefaultOptions.Authenticate(user.TfSecret, form.TfCode) {
		if models.SetTwoFactor(user, "", false) {
			models.SetMessage(user.Id, "关闭双重认证成功", "", "", "", "")
			u.Ok()
		} else {
			u.Error(100102)
		}
	} else {
		u.Error(100108)
	}
}

// @router /self/send-bind-captcha [get]
func (u *SelfController) SendBindCaptcha() {
	ident := u.GetString("ident")
	countryCode := u.GetString("country_code")

	user := u.GetUser()
	if strings.Contains(ident, "@") {
		if user.Email != "" {
			u.Error(100104)
			return
		}
		if !models.IsUsableEmail(ident) {
			u.Error(100305)
			return
		}
		if err := models.SendEmail(ident, "BIND"); err != nil {
			u.Error(err)
			return
		}
	} else {
		if user.Mobile != "" {
			u.Error(100104)
			return
		}
		if !models.IsUsableMobile(ident) {
			u.Error(100313)
			return
		}
		if err := models.SendSMS(ident, "BIND", countryCode); err != nil {
			u.Error(err)
			return
		}
	}
	u.Ok()
}

// @router /self/bind [post]
func (u *SelfController) Bind() {
	var form models.UserDto
	json.Unmarshal(u.Ctx.Input.RequestBody, &form)

	user := u.GetUser()
	captcha := utils.CaptchaGet(form.CountryCode + form.Ident, "BIND")
	if captcha == "" {
		u.Error(100307)
		return
	}
	if captcha != form.Captcha {
		u.Error(100308)
		return
	}

	if strings.Contains(form.Ident, "@") {
		if user.Email != "" {
			u.Error(100104)
			return
		}
		if !models.IsUsableEmail(form.Ident) {
			u.Error(100305)
			return
		}
		if err := models.BindEmail(user, form.Ident); err != nil {
			u.Error(err)
			return
		}
		message := "邮箱绑定成功"
		messageEn := "Successfully linked with your email."
		messageKo := "이메일이 성공적으로 연관되었습니다."
		messageJp := "メールアドレス設定成功"
		models.SetMessage(user.Id, message, messageEn, messageKo, messageJp, "")
	} else {
		if user.Mobile != "" {
			u.Error(100104)
			return
		}
		if !models.IsUsableMobile(form.Ident) {
			u.Error(100313)
			return
		}
		if err := models.BindMobile(user, form.Ident, form.CountryCode); err != nil {
			u.Error(err)
			return
		}
		message := "手机号绑定成功"
		messageEn := "Successfully linked with your phone."
		messageKo := "핸드폰 번호가 성공적으로 연관되었습니다."
		messageJp := "携帯番号設定成功"
		models.SetMessage(user.Id, message, messageEn, messageKo, messageJp, "")
	}

	u.Ok()
}

// @router /self/verify-code [get]
func (u *SelfController) SendVerifyCode() {
	mode := u.GetString("mode")
	captchaType := u.GetString("type")

	if mode != "EMAIL" && mode != "SMS" {
		u.Error("invalid mode")
		return
	}

	if captchaType != "TRANSFEROUT" && captchaType != "FUNDPASSWORD" && captchaType != "TRANSFERLOCKED" && captchaType != "TRANSFERUNLOCK" {
		u.Error("invalid type")
		return
	}

	user := u.GetUser()
	if mode == "EMAIL" {
		if user.Email == "" {
			u.Error(100320)
			return
		}
		if err := models.SendEmail(user.Email, captchaType); err != nil {
			u.Error(err)
			return
		}
	}
	if mode == "SMS" {
		if user.Mobile == "" {
			u.Error(100321)
			return
		}
		if err := models.SendSMS(user.Mobile, captchaType, user.CountryCode); err != nil {
			u.Error(err)
			return
		}
	}
	u.Ok()
}

// @router /self/fund-password [post]
func (u *SelfController) FundPassword() {
	type Form struct {
		VerifyMode string  `json:"verify_mode"`
		VerifyCode  string  `json:"verify_code"`
		FundPassword string `json:"fund_password"` //资金密码
		CheckPassword string `json:"check_password"` //确认密码
	}
	var form Form
	json.Unmarshal(u.Ctx.Input.RequestBody, &form)

	if form.FundPassword != form.CheckPassword {
		u.Error(100318)
		return
	}

	user := u.GetUser()
	switch form.VerifyMode {
	case "EMAIL":
		if user.Email == "" {
			u.Error(100320)
			return
		}
		captcha := utils.CaptchaGet(user.Email, "FUNDPASSWORD")
		if captcha == "" {
			u.Error(100307)
			return
		}
		if captcha != form.VerifyCode {
			u.Error(100308)
			return
		}
	case "SMS":
		if user.Mobile == "" {
			u.Error(100321)
			return
		}
		captcha := utils.CaptchaGet(user.CountryCode + user.Mobile, "FUNDPASSWORD")
		if captcha == "" {
			u.Error(100307)
			return
		}
		if captcha != form.VerifyCode {
			u.Error(100308)
			return
		}
	default:
		u.Error("invalid verify_mode")
		return
	}

	if err := models.SetFundPassword(user, form.FundPassword); err != nil {
		u.Error(err)
		return
	}
	if form.VerifyMode == "EMAIL" {
		utils.CaptchaDel(user.Email, "FUNDPASSWORD")
	}
	if form.VerifyMode == "SMS" {
		utils.CaptchaDel(user.CompleteMobile(), "FUNDPASSWORD")
	}

	message := "设置资金密码成功"
	messageEn := "Successfully set the fund password."
	messageKo := "자금 비밀번호가 성공적으로 설정되었습니다."
	messageJp := "資金パスワード設定成功"
	models.SetMessage(user.Id, message, messageEn, messageKo, messageJp, "")
	u.Ok()
}

// @router /self/messages [get]
func (u *SelfController) Messages() {
	page, err := u.GetInt64("page", 1)
	if err != nil {
		page = 1
	}
	user := u.GetUser()
	u.Data["json"] = result.Success(models.GetMessages(user, page))
	u.ServeJSON()
}

// @router /self/messages [post]
func (u *SelfController) MessagesAllRead() {
	user := u.GetUser()
	models.MessagesAllRead(user)
	u.Ok()
}

// @router /self/msg-num [get]
func (u *SelfController) GetMsgNum() {
	user := u.GetUser()
	u.Data["json"] = result.Success(models.GetMsgNum(user))
	u.ServeJSON()
}

// @router /self/message/:id [get]
func (u *SelfController) MessageRead() {
	id, err := u.GetInt(":id")
	if err != nil {
		u.Error("invalid id")
		return
	}
	user := u.GetUser()
	models.MessageRead(user, id)
	u.Ok()
}

// @router /self/fishing-code [post]
func (u *SelfController) SetFishingCode() {
	type Form struct {
		FishingCode string `json:"fishing_code"`
	}
	var form Form
	json.Unmarshal(u.Ctx.Input.RequestBody, &form)

	user := u.GetUser()
	if err := models.SetFishingCode(user, form.FishingCode); err!= nil {
		u.Error(err)
		return
	}
	models.SetMessage(user.Id, "设置防钓鱼码成功", "", "", "","")
	u.Ok()
}

// @router /self/fishing-code [get]
func (u *SelfController) GetFishingCode() {
	user := u.GetUser()
	utils.FishingCode.WriteTo(user.FishingCode, u.Ctx.ResponseWriter)
}

// @router /self/invitelink [get]
func (u *SelfController) GetInvitelink() {
	user := u.GetUser()
	u.Data["json"] = result.Success(models.GetInvitelink(user))
	u.ServeJSON()
}

// @router /self/invite-num [get]
func (u *SelfController) GetInviteNum() {
	user := u.GetUser()
	u.Data["json"] = result.Success(models.GetInviteNum(user))
	u.ServeJSON()
}

// @router /self/invites [get]
func (u *SelfController) GetInvites() {
	username := u.GetString("username")
	user := u.GetUser()
	u.Data["json"] = result.Success(models.GetInvites(user, username))
	u.ServeJSON()
}
