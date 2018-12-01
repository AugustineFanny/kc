package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"kuangchi_backend/models"
	"kuangchi_backend/result"
	"kuangchi_backend/utils"
	"strconv"
	"time"
)

type AdminController struct {
	CommonController
}

// @router /admin-login [post]
func (u *AdminController) Login() {
	var userDto models.UserDto
	json.Unmarshal(u.Ctx.Input.RequestBody, &userDto)

	if userDto.Code != u.GetSession("code") {
		u.Error(100101)
		return
	}

	admin, err := models.AdminSignIn(userDto.Ident, userDto.Password)
	if err != nil {
		u.Error(err)
		return
	}

	u.SetSession("admin", admin)
	u.Data["json"] = result.Success(admin)
	models.LogOperation(admin , "login", "", "", admin.Id,true)
	u.ServeJSON()
}

// @router /admin/logout [get]
func (u *AdminController) Logout() {
	u.DelSession("admin")
	u.Ok()
}

// @router /admin/change-password [post]
func (u *AdminController) ChangePassword() {
	admin := u.GetAdmin()

	var dto models.UserDto
	json.Unmarshal(u.Ctx.Input.RequestBody, &dto)

	err := models.ChangeAdminPassword(admin, dto.Password, dto.NewPassword)
	if err != nil {
		u.Error(err)
		models.LogOperation(admin, "change-password", nil, "Admin", admin.Id,false, err)
		return
	}
	models.LogOperation(admin , "change-password", nil, "Admin", admin.Id,true)
	u.Ok()
}

// @router /admin/set-mobile [post]
func (u *AdminController) SetMobile() {
	var form struct {
		Mobile    string  `json:"mobile"`
	}

	json.Unmarshal(u.Ctx.Input.RequestBody, &form)

	admin := u.GetAdmin()
	admin.Mobile = form.Mobile
	err := models.SetAdminMobile(*admin)
	if err != nil {
		u.Error(err)
		return
	}
	u.Ok()
}

// @router /admin/users [get]
func (u *AdminController) Users() {
	page, err := u.GetInt64("page", 1)
	if err != nil {
		page = 1
	}
	sort := u.GetString("_sort")
	order := u.GetString("_order")
	ident := u.GetString("ident")
	groupId, _ := u.GetInt("group_id", -1)
	u.Data["json"] = result.Success(models.GetUsers(page, sort, order, ident, groupId))
	u.ServeJSON()
}

// @router /admin/real-names [get]
func (u *AdminController) RealNames() {
	page, err := u.GetInt64("page", 1)
	if err != nil {
		page = 1
	}
	uid, _ := u.GetInt64("uid", -1)
	u.Data["json"] = result.Success(models.GetRealNames(page, uid))
	u.ServeJSON()
}

// @router /admin/name-auth [post]
func (u *AdminController) NameAuth() {
	var status models.StatusDto
	json.Unmarshal(u.Ctx.Input.RequestBody, &status)
	if !(status.Status == 2 || status.Status == 3) {
		u.Error("无效的status")
		return
	}
	admin := u.GetAdmin()
	if err := models.ChangeRole(status); err != nil {
		beego.Error(err)
		u.Error(err)
		models.LogOperation(admin, "name-auth", status, "RealName", status.Id,false, err)
		return
	}
	models.LogOperation(admin, "name-auth", status, "RealName", status.Id,true)
	u.Ok()
}

// @router /admin/kycs [get]
func (u*AdminController) Kycs() {
	page, err := u.GetInt64("page", 1)
	if err != nil {
		page = 1
	}

	uid, _ := u.GetInt64("uid", -1)
	u.Data["json"] = result.Success(models.GetKycs(page, uid))
	u.ServeJSON()
}

// @router /admin/kyc-auth [post]
func (u *AdminController) KycAuth() {
	var status models.StatusDto
	json.Unmarshal(u.Ctx.Input.RequestBody, &status)
	if !(status.Status == 2 || status.Status == 3) {
		u.Error("无效的status")
		return
	}
	admin := u.GetAdmin()
	if err := models.ChangeKyc(status); err != nil {
		beego.Error(err)
		u.Error(err)
		models.LogOperation(admin, "kyc-auth", status, "Kyc", status.Id,false, err)
		return
	}
	models.LogOperation(admin, "kyc-auth", status, "Kyc", status.Id,true)
	u.Ok()
}

// @router /admin/super/currencies [get]
func (u *AdminController) Currencies() {
	u.Data["json"] = result.Success(models.GetCurrencies("all"))
	u.ServeJSON()
}

// @router /admin/super/currencies [post]
func (u *AdminController) CreateCurrency() {
	var form models.KcCurrency
	json.Unmarshal(u.Ctx.Input.RequestBody, &form)
	if form.MInterestRate < 0 || form.MInterestRate > 0.005 {
		u.Error("MInterestRate 不能大于0.005 小于0")
		return
	}
	if form.CInterestRate < 0 || form.CInterestRate > 0.005 {
		u.Error("CInterestRate 不能大于0.005 小于0")
		return
	}
	admin := u.GetAdmin()
	id, err := models.CreateCurrency(form)
	if err != nil {
		beego.Error(err)
		u.Error(err)
		models.LogOperation(admin, "currencies", form, "KcCurrencies", int(*id),false, err)
		return
	}
	models.LogOperation(admin, "currencies", form, "KcCurrencies", int(*id),true)
	u.Ok()
}

// @router /admin/transfers [get]
func (u *AdminController) Transfers() {
	page, err := u.GetInt64("page", 1)
	if err != nil {
		page = 1
	}
	uid, _ := u.GetInt("uid", -1)
	currency := u.GetString("currency")
	status, _ := u.GetInt("status", -1)
	direction, _ := u.GetInt("direction", -1)
	u.Data["json"] = result.Success(models.GetTransfers(uid, currency, status, direction, page))
	u.ServeJSON()
}

// @router /admin/transfer-auth [post]
func (u *AdminController) TransferAuth() {
	var status models.StatusDto
	json.Unmarshal(u.Ctx.Input.RequestBody, &status)
	if !(status.Status == 2 || status.Status == 4) {
		u.Error("无效的status")
		return
	}
	admin := u.GetAdmin()
	if err := models.TransferAuth(status); err != nil {
		beego.Error(err)
		u.Error(err)
		models.LogOperation(admin, "transfer-auth", status, "UserAddressRecord", status.Id,false, err)
		return
	}
	models.LogOperation(admin, "transfer-auth", status, "UserAddressRecord", status.Id,true)
	u.Ok()
}

// @router /admin/super/admins [get]
func (u *AdminController) Admins() {
	u.Data["json"] = result.Success(models.GetAdmins())
	u.ServeJSON()
}

// @router /admin/super/admins [post]
func (u *AdminController) CreateAdmin() {
	var form models.UserDto
	json.Unmarshal(u.Ctx.Input.RequestBody, &form)

	admin := u.GetAdmin()
	row := models.Admin{Username: form.Ident, Password: form.Password}
	id, err := models.CreateAdmin(&row)
	if err != nil {
		beego.Error(err)
		u.Error(err)
		models.LogOperation(admin, "admins", nil, "Admin", int(*id),false, err)
		return
	}
	models.LogOperation(admin , "admins", nil, "Admin", int(*id),true)
	u.Ok()
}

// @router /admin/super/admin/{id} [post]
func (u *AdminController) ChangeAdmin() {
	var form struct {
		Username string `json:"username"`
		Status   int `json:"status"`
	}

	json.Unmarshal(u.Ctx.Input.RequestBody, &form)

	model := models.Admin{Username:form.Username, Status:form.Status}
	var err error
	if model.Id, err = u.GetInt(":id"); err != nil {
		u.Error("无效id")
		return
	}

	admin := u.GetAdmin()
	if err := models.SetAdmin(model); err != nil {
		beego.Error(err)
		u.Error(err)
		models.LogOperation(admin, "admin/{id}", form, "Admin", model.Id,false, err)
		return
	}
	models.LogOperation(admin , "admin/{id}", form, "Admin", model.Id,true)
	u.Ok()
}

// @router /admin/super/currency/:id [post]
func (u *AdminController) ChangeCurrency() {
	var form models.KcCurrency
	json.Unmarshal(u.Ctx.Input.RequestBody, &form)

	var err error
	if form.Id, err = u.GetInt(":id"); err != nil {
		u.Error("无效id")
		return
	}
	if form.MInterestRate < 0 || form.MInterestRate > 0.005 {
		u.Error("MInterestRate 不能大于0.005 小于0")
		return
	}
	if form.CInterestRate < 0 || form.CInterestRate > 0.005 {
		u.Error("CInterestRate 不能大于0.005 小于0")
		return
	}
	admin := u.GetAdmin()
	if err := models.SetCurrency(form); err != nil {
		beego.Error(err)
		u.Error(err)
		models.LogOperation(admin, "currency/{id}", form, "Currency", form.Id,false, err)
		return
	}
	models.LogOperation(admin , "currency/{id}", form, "Currency", form.Id,true)
	u.Ok()
}

// @router /admin/super/logs [get]
func (u *AdminController) Logs() {
	page, err := u.GetInt64("page", 1)
	if err != nil {
		page = 1
	}
	ident := u.GetString("ident")
	api := u.GetString("api")
	status, _ := u.GetInt("status", -1)
	u.Data["json"] = result.Success(models.GetLogs(page, ident, api, status))
	u.ServeJSON()
}

// @router /admin/user/thaw [post]
func (u *AdminController) Thaw() {
	var form struct {
		Username     string
	}
	json.Unmarshal(u.Ctx.Input.RequestBody, &form)
	admin := u.GetAdmin()
	utils.RedisClient.Delete("blocked:" + form.Username)
	models.LogOperation(admin , "blocked", form, "解冻账户", 0,true)
	u.Ok()
}

// @router /admin/wallets [get]
func (u *AdminController) Wallets() {
	page, err := u.GetInt64("page", 1)
	if err != nil {
		page = 1
	}
	ident := u.GetString("ident")
	currency := u.GetString("currency")
	noZero, _ := u.GetBool("no_zero", false)
	sort := u.GetString("_sort")
	order := u.GetString("_order")
	u.Data["json"] = result.Success(models.BackGetWallets(page, ident, currency, noZero, sort, order))
	u.ServeJSON()
}

// @router /admin/addresses [get]
func (u *AdminController) Addresses() {
	page, err := u.GetInt64("page", 1)
	if err != nil {
		page = 1
	}
	uid, _ := u.GetInt("uid", -1)
	currency := u.GetString("currency")
	noZero, _ := u.GetBool("no_zero", false)
	u.Data["json"] = result.Success(models.BackGetAddresses(page, uid, currency, noZero))
	u.ServeJSON()
}

// @router /admin/address/pool [get]
func (u *AdminController) AddressPool() {
	page, err := u.GetInt64("page", 1)
	if err != nil {
		page = 1
	}
	currency := u.GetString("currency")
	flag, _ := u.GetInt("flag", -1)
	u.Data["json"] = result.Success(models.BackGetAddressePool(page, currency, flag))
	u.ServeJSON()
}

// @router /admin/statistics [get]
func (u *AdminController) Statistics() {
	res := map[string]interface{}{
		"balance": models.Stat("balance"),
		"recharge": models.Stat("recharge"),
		"withdraw": models.Stat("withdraw"),
		"sign_up_today": models.Stat("sign_up_today"),
		"sign_up_yesterday": models.Stat("sign_up_yesterday"),
		"sign_up_month": models.Stat("sign_up_month"),
		"sign_up_prev_month": models.Stat("sign_up_prev_month"),
		"fet_stat": models.Stat("fet_stat"),
	}
	u.Data["json"] = result.Success(res)
	u.ServeJSON()
}

// @router /admin/invites [get]
func (u *AdminController) Invites() {
	startDate := u.GetString("start_date")
	endDate := u.GetString("end_date")
	u.Data["json"] = result.Success(models.BackGetInvites(startDate, endDate))
	u.ServeJSON()
}

// @router /admin/child/amounts [get]
func (u *AdminController) ChildAmounts() {
	u.Data["json"] = result.Success(models.BackGetChildAmounts())
	u.ServeJSON()
}

// @router /admin/groups [get]
func (u *AdminController) Groups() {
	u.Data["json"] = result.Success(models.GetGroups())
	u.ServeJSON()
}

// @router /admin/super/groups [post]
func (u *AdminController) CreateGroup() {
	var form struct {
		Name      string   `json:"name"`
	}
	json.Unmarshal(u.Ctx.Input.RequestBody, &form)
	admin := u.GetAdmin()
	id, err := models.CreateGroup(form.Name)
	if err != nil {
		beego.Error(err)
		u.Error(err)
		models.LogOperation(admin, "groups", form, "Group", int(*id),false, err)
		return
	}
	models.LogOperation(admin, "groups", form, "Group", int(*id),true)
	u.Ok()
}

// @router /admin/super/user/:id [post]
func (u *AdminController) ChangeUser() {
	var form struct {
		GroupId    int   `json:"group_id"`
	}
	json.Unmarshal(u.Ctx.Input.RequestBody, &form)

	model := models.User{GroupId: form.GroupId}
	var err error
	if model.Id, err = u.GetInt(":id"); err != nil {
		u.Error("无效id")
		return
	}

	admin := u.GetAdmin()
	if err := models.ChangeUser(model); err != nil {
		beego.Error(err)
		u.Error(err)
		models.LogOperation(admin, "user/{id}", form, "User", model.Id,false, err)
		return
	}
	models.LogOperation(admin, "user/{id}", form, "User", model.Id,true)
	u.Ok()
}

// @router /admin/super/captcha [get]
func (u *AdminController) SendCaptcha() {
	captchaType := u.GetString("type")
	if err := models.BackSendCaptcha(utils.SuperMobile, captchaType); err != nil {
		u.Error(err)
		return
	}
	u.Ok()
}

// @router /admin/super/super-password [post]
func (u *AdminController) SuperPassword() {
	var form models.UserDto
	json.Unmarshal(u.Ctx.Input.RequestBody, &form)
	if form.Password != form.CheckPassword {
		u.Error("两次密码不一致")
		return
	}
	if !utils.RedisClient.IsExist("captcha:superpassword:" + utils.SuperMobile) {
		u.Error("重新获取验证码")
		return
	}
	if form.Captcha == "" {
		u.Error("缺少参数 captcha")
		return
	}
	captcha, err := utils.RedisClient.GetString("captcha:superpassword:" + utils.SuperMobile)
	if err != nil {
		u.Error(err)
		return
	}
	if form.Captcha != captcha {
		u.Error(100308)
		return
	}
	admin := u.GetAdmin()
	err = models.SetSuperPassword(admin, form.Password)
	if err != nil {
		u.Error(err)
		models.LogOperation(admin, "super-password", nil, "Admin", admin.Id,false, err)
		return
	}
	utils.RedisClient.Delete("captcha:superpassword:" + utils.SuperMobile)
	models.LogOperation(admin , "super-password", nil, "Admin", admin.Id,true)
	u.Ok()
}

// @router /admin/super/user/:id/recharge [post]
func (u *AdminController) UserRecharge() {
	var form struct {
		Currency string `json:"currency"`
		Amount   float64 `json:"amount"`
		Password string  `json:"password"`
		Captcha  string  `json:"captcha"`
	}
	json.Unmarshal(u.Ctx.Input.RequestBody, &form)

	uid, err := u.GetInt(":id")
	if err != nil {
		u.Error("无效id")
		return
	}
	if form.Password == "" {
		u.Error("缺少参数 password")
		return
	}
	if form.Amount <= 0 {
		u.Error("无效参数 amount")
		return
	}
	if form.Currency == "" || form.Currency == "BTC" || form.Currency == "ETH" || form.Currency == "USDT" {
		u.Error("无效币种")
		return
	}

	admin := u.GetAdmin()
	if !admin.ValidateFundPassword(form.Password) {
		u.Error("密码错误")
		return
	}
	idstr := strconv.Itoa(admin.Id)
	if !utils.RedisClient.IsExist("captcha:userrecharge:" + idstr) {
		if form.Captcha == "" {
			u.Error(900001)
			return
		} else {
			captcha, err := utils.RedisClient.GetString("captcha:userrecharge:" + utils.SuperMobile)
			if err != nil {
				u.Error(err)
				return
			}
			if form.Captcha != captcha {
				u.Error(100308)
				return
			}
			utils.RedisClient.Set("captcha:userrecharge:" + idstr, "true", 60*60)
		}
	}
	if err := models.UserRecharge(uid, form.Currency, form.Amount); err != nil {
		beego.Error(err)
		u.Error(err)
		form.Password = ""
		form.Captcha = ""
		models.LogOperation(admin, "recharge/{id}", form, "Wallet", 0,false, err)
		return
	}
	form.Password = ""
	form.Captcha = ""
	models.LogOperation(admin, "recharge/{id}", form, "Wallet", 0,true)
	u.Ok()
}

// @router /admin/super/fund-changes [get]
func (u *AdminController) FundChanges() {
	page, err := u.GetInt64("page", 1)
	if err != nil {
		page = 1
	}
	uid, _ := u.GetInt("uid", -1)
	currency := u.GetString("currency")
	direction, _ := u.GetInt("direction", -1)
	desc := u.GetString("desc")
	u.Data["json"] = result.Success(models.BackGetFundChanges(page, uid, currency, direction, desc))
	u.ServeJSON()
}

// @router /admin/super/batch-sms [post]
func (u *AdminController) BatchSms() {
	var form struct {
		Batch       bool       `json:"batch"`
		Mobile      string     `json:"mobile"`
		Content     string     `json:"content"`
	}
	json.Unmarshal(u.Ctx.Input.RequestBody, &form)

	if form.Batch {
		form.Mobile = ""
	} else if form.Mobile == "" {
		u.Error("参数不能为空 mobile")
		return
	}

	admin := u.GetAdmin()
	if err := models.BatchSms(form.Content, form.Mobile); err != nil {
		beego.Error(err)
		u.Error(err)
		models.LogOperation(admin, "batch-sms", form, "", 0,false, err)
		return
	}
	models.LogOperation(admin , "batch-sms", form, "", 0,true)
	u.Ok()
}

// @router /admin/super/batch-email [post]
func (u *AdminController) BatchEmail() {
	var form struct {
		Batch       bool       `json:"batch"`
		Email       string     `json:"email"`
		Content     string     `json:"content"`
	}
	json.Unmarshal(u.Ctx.Input.RequestBody, &form)

	if form.Batch {
		form.Email = ""
	} else {
		if form.Email == "" {
			u.Error("参数不能为空 email")
			return
		}
		valid := validation.Validation{}
		if v := valid.Email(form.Email, "email"); !v.Ok {
			u.Error("邮箱格式不正确")
			return
		}
	}

	admin := u.GetAdmin()
	if err := models.BatchEmail(form.Content, form.Email); err != nil {
		beego.Error(err)
		u.Error(err)
		models.LogOperation(admin, "batch-email", form, "", 0,false, err)
		return
	}
	models.LogOperation(admin , "batch-email", form, "", 0,true)
	u.Ok()
}

// @router /admin/super/captcha-to-user [get]
func (u *AdminController) BackSendCaptchaToUser() {
	username := u.GetString("username")
	if err := models.SendCaptchaToUserByUsername(username); err != nil {
		u.Error(err)
		return
	}
	u.Ok()
}

// @router /admin/profit-date [get]
func (u *AdminController) ProfitDate() {
	desc := u.GetString("desc", "mining")
	date := u.GetString("date", time.Now().Format("2006-01-02"))
	u.Data["json"] = result.Success(models.GetProfitForDate(desc, date))
	u.ServeJSON()
}

// @router /admin/profit-month [get]
func (u *AdminController) ProfitMonth() {
	desc := u.GetString("desc", "mining")
	startDate := u.GetString("start_date")
	endDate := u.GetString("end_date")
	u.Data["json"] = result.Success(models.GetProfitForMonth(desc, startDate, endDate))
	u.ServeJSON()
}

// @router /admin/subscriptions [get]
func (u *AdminController) GetSubscriptions() {
	page, err := u.GetInt64("page", 1)
	if err != nil {
		page = 1
	}
	ident := u.GetString("ident")
	order := u.GetString("order")
	status, _ := u.GetInt("status", -1)
	base := u.GetString("base")
	date := u.GetString("date")
	u.Data["json"] = result.Success(models.BackGetSubscriptions(page, ident, order, base, status, date))
	u.ServeJSON()
}

// @router /admin/order/:order/submissions [get]
func (u *AdminController) GetSubmissions() {
	order := u.GetString(":order")
	submissions, err := models.BackGetSubmissions(order)
	if err != nil {
		u.Error(err)
		return
	}
	u.Data["json"] = result.Success(submissions)
	u.ServeJSON()
}

// @router /admin/order [get]
func (u *AdminController) ConfirmOrder() {
	var args struct {
		Id        int
		Status    int
		Remark    string
	}

	json.Unmarshal(u.Ctx.Input.RequestBody, &args)
	if !(args.Status == 0 || args.Status == 2) {
		u.Error("无效的status")
		return
	}
	admin := u.GetAdmin()
	if err := models.ConfirmOrder(args.Id, args.Status, args.Remark); err != nil {
		beego.Error(err)
		u.Error(err)
		models.LogOperation(admin, "order", args, "kc_subscription", args.Id,false, err)
		return
	}
	models.LogOperation(admin, "order", args, "kc_subscription", args.Id,true)
	u.Ok()
}
