package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"kuangchi_backend/utils"
)
func GetIUUExchange() map[string]float64 {
	var curs []*KcCurrency
	res := map[string]float64{
		"USDT": 0,
		"BTC": 0,
		"ETH": 0,
	}
	o := orm.NewOrm()
	if _, err := o.QueryTable("kc_currency").Filter("currency__in", "USDT", "BTC", "ETH", "IUU").All(&curs); err != nil {
		beego.Error(err)
		return res
	}
	temp := map[string]float64{}
	for _, cur := range curs {
		temp[cur.Currency] = cur.BasePrice
	}
	if _, ok := temp["IUU"]; !ok {
		beego.Error("not found IUU")
		return res
	}
	for cur, basePrice := range temp {
		if cur != "IUU" {
			res[cur] = utils.ShowFloat(basePrice/temp["IUU"], 2)
		}
	}
	return res
}