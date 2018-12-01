package controllers

import (
	"kuangchi_backend/result"
	"kuangchi_backend/utils"
	"kuangchi_backend/models"
)

type PublicController struct {
	CommonController
}

// @router /public/currencies [get]
func (u *PublicController) Currencies() {
	flag := u.GetString("flag")
	u.Data["json"] = result.Success(models.GetCurrencies(flag))
	u.ServeJSON()
}

// @router /public/currencies/detail [get]
func (u *PublicController) CurrenciesDetail() {
	flag := u.GetString("flag")
	u.Data["json"] = result.Success(models.GetCurrencies(flag))
	u.ServeJSON()
}

// @router /public/payment-methods [get]
func (u *PublicController) PaymentMethods() {
	u.Data["json"] = result.Success(utils.PaymentMethods)
	u.ServeJSON()
}

// @router /public/units [get]
func (u *PublicController) Units() {
	u.Data["json"] = result.Success(utils.Units)
	u.ServeJSON()
}

// @router /public/countries [get]
func (u *PublicController) Countries() {
	u.Data["json"] = result.Success(utils.Countries)
	u.ServeJSON()
}

// @router /public/exchanges [get]
func (u *PublicController) Exchanges() {
	currency := u.GetString("currency")
	u.Data["json"] = result.Success(models.GetValidExchange(currency))
	u.ServeJSON()
}

// @router /public/exchange/:exchange [get]
func (u *PublicController) Exchange() {
	exchange := u.GetString(":exchange")
	currency := u.GetString("currency")
	unit := u.GetString("unit")
	premium, err := u.GetFloat("premium", 0)
	if err != nil {
		premium = 0
	}
	u.Data["json"] = result.Success(utils.RealTimePrice(exchange, currency, unit, premium))
	u.ServeJSON()
}

// @router /public/:currency/nodes [get]
func (u *PublicController) Nodes() {
	currency := u.GetString(":currency")
	if !models.ValidCurrency(currency) {
		u.Error(100401)
		return
	}
	u.Data["json"] = result.Success(models.GetNodes(currency))
	u.ServeJSON()
}

// @router /public/IUU/exchange [get]
func (u *PublicController) IUUExchange() {
	u.Data["json"] = result.Success(models.GetIUUExchange())
	u.ServeJSON()
}
