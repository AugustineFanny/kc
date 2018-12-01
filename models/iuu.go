package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"kuangchi_backend/result"
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

func CreateOrder(u *User, dest, base string, amount, curAmount float64) error {
	exchanges := GetIUUExchange()
	exchange, ok := exchanges[base]
	if !ok {
		return result.ErrCode(100102)
	}
	curAmount = utils.ShowFloat(amount * exchange, 2)
	if exchange > 0 && curAmount < 0.01 {
		return result.ErrCode(100406)
	}
	o := orm.NewOrm()
	o.Begin()
	if err := ReduceAmount(o, u.Id, amount, false, base, "subscription"); err != nil {
		o.Rollback()
		return err
	}
	if err := AddAmount(o, u.Id, curAmount, dest, "subscription"); err != nil {
		o.Rollback()
		beego.Error(err)
		return result.ErrCode(100102)
	}
	//认购直接锁仓
	if err := handleLocked(o, u, dest, curAmount * 0.2, 0); err != nil {
		o.Rollback()
		return err
	}
	subscription := KcSubscription{
		Uid: u.Id,
		Currency: dest,
		Base: base,
		BaseAmount:amount,
		CurAmount:curAmount,
		Exchange: exchange,
	}
	if _, err := o.Insert(&subscription); err != nil {
		o.Rollback()
		beego.Error(err)
		return result.ErrCode(100102)
	}
	o.Commit()

	message := fmt.Sprintf("获得 %f %s", curAmount, dest)
	messageEn := fmt.Sprintf("Get %f %s", curAmount, dest)
	messageKo := fmt.Sprintf("얻다 %f %s", curAmount, dest)
	messageJp := fmt.Sprintf("獲得 %f %s", curAmount, dest)
	SetMessage(u.Id, message, messageEn, messageKo, messageJp, "")
	return nil
}