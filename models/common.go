package models

import (
	"kuangchi_backend/utils"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"time"
)

func SendLoginNotifination(u *User) {
	//if u.Verified("mobile") {
	//	go utils.PushSms("This is to notify you of a successful login to your account. Login Time: " + time.Now().Format("2006-01-02 15:04:05"), u.Mobile)
	//}
	if u.Verified("email") {
		go utils.SendContentEmail(u.Email, "FADAX", "This is to notify you of a successful login to your account. Login Time: " + time.Now().Format("2006-01-02 15:04:05"))
	}
}

func SendNotifination(u *User, content string) {
	if u != nil {
		//if u.Verified("mobile") {
		//	go utils.PushSms(content + "【投币时代】", u.Mobile)
		//}
		if u.Verified("email") {
			go utils.SendContentEmail(u.Email, "FADAX", content)
		}
	}
}

func GetCurrencies(flag string) []*KcCurrency {
	var err error
	o := orm.NewOrm()
	var currencies []*KcCurrency
	if flag == "all" {
		_, err = o.QueryTable("kc_currency").All(&currencies)
	} else {
		_, err = o.QueryTable("kc_currency").Filter("trans_flag", 0).All(&currencies)
	}
	if err != nil {
		beego.Error(err)
		return []*KcCurrency{}
	}
	return currencies
}

func GetCurrencyStrings(flag string) []string {
	var err error
	o := orm.NewOrm()
	var list orm.ParamsList
	if flag == "all" {
		_, err = o.QueryTable("kc_currency").ValuesFlat(&list, "currency")
	} else {
		_, err = o.QueryTable("kc_currency").Filter("trans_flag", 0).ValuesFlat(&list, "currency")
	}
	if err != nil {
		beego.Error(err)
		return []string{}
	}
	newList := make([]string, len(list))
	for i, v := range list {
		newList[i] = v.(string)
	}
	return newList
}

func GetCurrency(currency string) *KcCurrency {
	o := orm.NewOrm()
	cur := KcCurrency{Currency:currency}
	if err := o.Read(&cur, "Currency"); err != nil {
		return nil
	}
	return &cur
}

func GetValidExchange(currency string) []string {
	cur := GetCurrency(currency)
	if cur != nil {
		return cur.ExchangesList()
	}
	return []string {}
}

func ValidCurrency(currency string) bool {
	o := orm.NewOrm()
	if o.QueryTable("kc_currency").Filter("currency", currency).Exist() {
		return true
	}
	return false
}

func ValidExchange(currency, exchange string) bool {
	cur := GetCurrency(currency)
	if cur != nil {
		for _, e := range cur.ExchangesList() {
			if e == exchange {
				return true
			}
		}
	}
	return false
}

func TrustRelate(uid, fid int) bool {
	o := orm.NewOrm()
	return o.QueryTable("trust").Filter("UserId", uid).Filter("FollowId", fid).Exist()
}
