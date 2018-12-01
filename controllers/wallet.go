package controllers

import (
	"encoding/json"
	"kuangchi_backend/models"
	"kuangchi_backend/result"
	"kuangchi_backend/utils"
	"fmt"
)

type WalletController struct {
	CommonController
}

// @router /wallet/currency/:currency [get]
func (u *WalletController) Address() {
	currency := u.GetString(":currency")
	if !models.ValidCurrency(currency) {
		u.Error(100401)
		return
	}
	user := u.GetUser()
	address := models.GetAddress(user, currency)
	data := map[string]string{"address": address}
	u.Data["json"] = result.Success(data)
	u.ServeJSON()
}

// @router /wallet/finance [get]
func (u *WalletController) Finance() {
	currency := u.GetString("currency")
	user := u.GetUser()
	if currency != "" {
		if !models.ValidCurrency(currency) {
			u.Error(100401)
			return
		}
		detail := models.GetWallet(user, currency)
		u.Data["json"] = result.Success(detail)
		u.ServeJSON()
		return
	}
	detail, err := models.GetFinance(user)
	if err != nil {
		u.Error(err)
		return
	}
	u.Data["json"] = result.Success(detail)
	u.ServeJSON()
}

// @router /wallet/transfer-out [post]
func (u *WalletController) TransferOut() {
	var form struct {
		Currency    string  `json:"currency"`
		Address     string  `json:"address"`
		Amount      float64 `json:"amount"`
		VerifyMode string  `json:"verify_mode"`
		VerifyCode  string  `json:"verify_code"`
		FundPassword string `json:"fund_password"`
		Method      string  `json:"method"` //withdraw:提现 instation:站内划转
	}
	json.Unmarshal(u.Ctx.Input.RequestBody, &form)
	if !models.ValidCurrency(form.Currency) {
		u.Error(100401)
		return
	}
	if form.Amount <= 0 {
		u.Error(100106)
		return
	}
	user := u.GetUser()
	//if user.Role != 3 {
	//	u.Error(100113)
	//	return
	//}
	if utils.RedisClient.IsExist("blocked:" + user.Username) {
		u.Error("Account freeze")
		return
	}
	switch form.VerifyMode {
	case "EMAIL":
		if user.Email == "" {
			u.Error(100320)
			return
		}
		captcha := utils.CaptchaGet(user.Email, "TRANSFEROUT")
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
		captcha := utils.CaptchaGet(user.CountryCode + user.Mobile, "TRANSFEROUT")
		if captcha == "" {
			u.Error(100307)
			return
		}
		if captcha != form.VerifyCode {
			u.Error(100308)
			return
		}
	case "TWOFACTOR":
		if user.TfOpened == false {
			u.Error("未开启两步验证")
			return
		}
		if !utils.DefaultOptions.Authenticate(user.TfSecret, form.VerifyCode) {
			u.Error(100108)
			return
		}
	default:
		u.Error("invalid verify_mode")
		return
	}
	if !user.ValidateFundPassword(form.FundPassword) {
		u.Error(100110)
		return
	}
	cur := models.GetCurrency(form.Currency)
	if cur == nil {
		u.Error(100102)
		return
	}
	if cur.FeeEth == 1 {
		//使用ETH做手续费
		if models.UsableBalance(user, "ETH") < cur.Fee {
			u.Error(100404)
			return
		}
	} else {
		//使用当前币做手续费
		if form.Amount <= cur.Fee {
			u.Error(100404)
			return
		}
	}
	if form.Method == "instation" {
		if err := models.TransferInStation(user, form.Currency, form.Address, form.Amount); err != nil {
			u.Error(err)
			return
		}
	} else {
		if err := models.ApplyWithdraw(user, form.Currency, form.Address, form.Amount, cur.Fee, cur.FeeEth); err != nil {
			u.Error(err)
			return
		}
	}
	if form.VerifyMode == "EMAIL" {
		utils.CaptchaDel(user.Email, "TRANSFEROUT")
	}
	if form.VerifyMode == "SMS" {
		utils.CaptchaDel(user.CompleteMobile(), "TRANSFEROUT")
	}
	u.Ok()
}

// @router /wallet/transfers [get]
func (u *WalletController) Transfers() {
	page, err := u.GetInt64("page", 1)
	if err != nil {
		page = 1
	}
	pageSize, err := u.GetInt64("page_size", 5)
	if err != nil {
		pageSize = 5
	}
	user := u.GetUser()
	currency := u.GetString("currency")
	direction, _ := u.GetInt("direction", -1)
	detail := models.GetTransfers(user.Id, currency, -1, direction, page, pageSize)
	if err != nil {
		u.Error(err)
		return
	}
	u.Data["json"] = result.Success(detail)
	u.ServeJSON()
}

// @router /wallet/usable/:currency [get]
func (u *WalletController) Usable() {
	currency := u.GetString(":currency")
	mode := u.GetString("mode")
	if !models.ValidCurrency(currency) {
		u.Error(100401)
		return
	}

	user := u.GetUser()
	rate := models.GetRate(user.Id, currency)
	var balance float64
	if mode == "TRADE" {
		balance = models.UsableBalance(user, currency, rate)
	} else if mode == "ALL" {
		balance = models.UsableBalance(user, currency)
	} else {
		u.Error("invalid mode")
		return
	}
	u.Data["json"] = result.Success(balance)
	u.ServeJSON()
}

// @router /wallet/:currency/locked [post]
func (u *WalletController) Locked() {
	currency := u.GetString(":currency")
	//if !models.ValidCurrency(currency) {
	//	u.Error(100401)
	//	return
	//}

	//只有FET能锁仓
	currency = "FET"
	var form struct {
		Amount    float64  `json:"amount"`
		FundPassword string `json:"fund_password"` //资金密码
	}
	json.Unmarshal(u.Ctx.Input.RequestBody, &form)
	if form.Amount <= 0 {
		u.Error("invalid amount")
		return
	}
	cur := models.GetCurrency(currency)
	if form.Amount < cur.MinLock {
		u.Error(fmt.Sprintf("min %f %s", cur.MinLock, currency))
		return
	}
	user := u.GetUser()
	if !user.ValidateFundPassword(form.FundPassword) {
		u.Error(100110)
		return
	}
	if err := models.Locked(user, currency, form.Amount); err != nil {
		u.Error(err)
		return
	}

	u.Ok()
}

// @router /wallet/:currency/locked [get]
func (u *WalletController) GetLocked() {
	currency := u.GetString(":currency")
	if !models.ValidCurrency(currency) {
		u.Error(100401)
		return
	}
	user := u.GetUser()
	u.Data["json"] = result.Success(models.GetLocked(user, currency))
	u.ServeJSON()
}

// @router /wallet/:currency/mining [get]
func (u *WalletController) GetMining() {
	currency := u.GetString(":currency")
	if !models.ValidCurrency(currency) {
		u.Error(100401)
		return
	}
	page, err := u.GetInt64("page", 1)
	if err != nil {
		page = 1
	}
	user := u.GetUser()
	u.Data["json"] = result.Success(models.GetMining(user, currency, page))
	u.ServeJSON()
}

// @router /wallet/:currency/mining-stat [get]
func (u *WalletController) GetMiningStat() {
	currency := u.GetString(":currency")
	if !models.ValidCurrency(currency) {
		u.Error(100401)
		return
	}
	user := u.GetUser()
	u.Data["json"] = result.Success(models.GetMiningStat(user, currency))
	u.ServeJSON()
}

// @router /wallet/locked/:id/transfer [post]
func (u *WalletController) TransferLocked() {
	id, err := u.GetInt(":id")
	if err != nil {
		u.Error("invalid")
		return
	}
	var form struct {
		Address     string  `json:"address"`
		VerifyMode string  `json:"verify_mode"`
		VerifyCode  string  `json:"verify_code"`
		FundPassword string `json:"fund_password"`
	}
	if err := json.Unmarshal(u.Ctx.Input.RequestBody, &form); err != nil {
		u.Error("invalid args")
		return
	}
	user := u.GetUser()
	switch form.VerifyMode {
	case "EMAIL":
		if user.Email == "" {
			u.Error(100320)
			return
		}
		captcha := utils.CaptchaGet(user.Email, "TRANSFERLOCKED")
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
		captcha := utils.CaptchaGet(user.CountryCode + user.Mobile, "TRANSFERLOCKED")
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
	if !user.ValidateFundPassword(form.FundPassword) {
		u.Error(100110)
		return
	}
	if err := models.TransferLocked(user, id, form.Address); err != nil {
		u.Error(err)
		return
	}
	if form.VerifyMode == "EMAIL" {
		utils.CaptchaDel(user.Email, "TRANSFERLOCKED")
	}
	if form.VerifyMode == "SMS" {
		utils.CaptchaDel(user.CompleteMobile(), "TRANSFERLOCKED")
	}
	u.Ok()
}

// @router /wallet/instations [get]
func (u *WalletController) Instations() {
	page, err := u.GetInt64("page", 1)
	if err != nil {
		page = 1
	}
	pageSize, err := u.GetInt64("page_size", 5)
	if err != nil {
		pageSize = 5
	}
	user := u.GetUser()
	currency := u.GetString("currency")
	direction, _ := u.GetInt("direction", -1)
	detail := models.GetInstations(user.Username, currency, direction, models.Instations, page, pageSize)
	if err != nil {
		u.Error(err)
		return
	}
	u.Data["json"] = result.Success(detail)
	u.ServeJSON()
}

// @router /wallet/inlocked [get]
func (u *WalletController) Inlocked() {
	page, err := u.GetInt64("page", 1)
	if err != nil {
		page = 1
	}
	pageSize, err := u.GetInt64("page_size", 5)
	if err != nil {
		pageSize = 5
	}
	user := u.GetUser()
	currency := u.GetString("currency")
	direction, _ := u.GetInt("direction", -1)
	detail := models.GetInstations(user.Username, currency, direction, models.Inlocked, page, pageSize)
	if err != nil {
		u.Error(err)
		return
	}
	u.Data["json"] = result.Success(detail)
	u.ServeJSON()
}

// @router /wallet/orders [post]
func (u *WalletController) CreateOrder() {
	var form struct {
		Dest         string   `json:"dest"`
		Base         string   `json:"base"`
		Amount       float64  `json:"amount"`
		CurAmount    float64  `json:"cur_amount"`
		FundPassword string   `json:"fund_password"`
	}
	json.Unmarshal(u.Ctx.Input.RequestBody, &form)
	if form.Amount <= 0 {
		u.Error("invalid amount")
		return
	}
	if form.Base != "USDT" && form.Base != "BTC" && form.Base != "ETH" {
		u.Error(100106)
		return
	}
	user := u.GetUser()
	if !user.ValidateFundPassword(form.FundPassword) {
		u.Error(100110)
		return
	}
	//目前只能认购IUU
	form.Dest = "IUU"
	err := models.CreateOrder(user, form.Dest, form.Base, form.Amount, form.CurAmount)
	if err != nil {
		u.Error(err)
		return
	}
	u.Ok()
}

// @router /wallet/orders [get]
func (u *WalletController) Orders() {
	page, err := u.GetInt64("page", 1)
	if err != nil {
		page = 1
	}
	user := u.GetUser()
	u.Data["json"] = result.Success(models.GetOrders(user, page))
	u.ServeJSON()
}

// @router /wallet/hashrate [get]
func (u *WalletController) Hashrate() {
	user := u.GetUser()
	u.Data["json"] = result.Success(models.GetHashrate(user))
	u.ServeJSON()
}
