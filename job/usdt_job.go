package job

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/httplib"
	"kuangchi_backend/models"
	"strings"
	"time"
	"github.com/spf13/cast"
)

type omniTransaction struct {
	Amount           string   `json:"amount"`
	Confirmations    int      `json:"confirmations"`
	Propertyid       int      `json:"propertyid"`
	ReferenceAddress string   `json:"referenceaddress"`
	SendingAddress   string   `json:"sendingaddress"`
	Txid             string   `json:"txid"`
	Valid            bool     `json:"valid"`
}

type omni struct {
	Address       string              `json:"address"`
	Transactions  []*omniTransaction  `json:"transactions"`
}

func usdtServer(address string) ([]*omniTransaction, error) {
	resp := omni{}
	url := "https://api.omniexplorer.info/v1/transaction/address/0"
	omniserver := httplib.Post(url).SetTimeout(10*time.Second, 10*time.Second)
	omniserver.Param("addr", address)
	if err := omniserver.ToJSON(&resp); err != nil {
		beego.Error(err)
		omniSleep()
		return nil, err
	}
	if resp.Address != address {
		return omni{}.Transactions, nil
	}
	return resp.Transactions, nil
}

func usdtDiposit(uid int, tran omniTransaction) {
	o := orm.NewOrm()
	o.Begin()
	record := models.KcUserAddressRecord{Uid: uid, Currency: "USDT", Hash: tran.Txid}
	amount := cast.ToFloat64(tran.Amount)
	created, _, err := o.ReadOrCreate(&record, "Uid", "Currency", "Hash")
	if err != nil {
		beego.Error(err)
		o.Rollback()
		return
	}
	if created {
		if tran.Valid == false {
			record.Status = 1
			record.CheckTime = time.Now()
		} else if tran.Confirmations > 1 {
			record.Status = 2
			record.CheckTime = time.Now()
			if err := models.AddAmount(o, uid, amount, "USDT", "deposit"); err != nil {
				beego.Error(err)
				o.Rollback()
				return
			}
		}
		record.Uid = uid
		record.From = tran.SendingAddress
		record.To = tran.ReferenceAddress
		record.Amount = amount
		record.Confirmations = tran.Confirmations
		if _, err := o.Update(&record); err != nil {
			beego.Error(err)
			o.Rollback()
			return
		}
	} else {
		if record.Status == 0 {
			if tran.Valid == false {
				record.Status = 1
				record.CheckTime = time.Now()
			} else if tran.Confirmations > 1 {
				record.Status = 2
				record.CheckTime = time.Now()
				if err := models.AddAmount(o, uid, amount, "USDT", "deposit"); err != nil {
					beego.Error(err)
					o.Rollback()
					return
				}
			}
			record.Confirmations = tran.Confirmations
			if _, err := o.Update(&record); err != nil {
				beego.Error(err)
				o.Rollback()
				return
			}
		}
	}

	o.Commit()
}

func omniSleep() {
	time.Sleep(time.Second * 10)
}

func USDTJob() {
	var lists []*models.KcUserAddress
	var offset int = 0
	o := orm.NewOrm()
	beego.Warn("usdt cron start")
	for {
		num, err := o.QueryTable("kc_user_address").Filter("currency", "USDT").Limit(100, offset).All(&lists, "uid", "address")
		if err != nil {
			beego.Error(err)
			time.Sleep(time.Second * 10)
			continue
		}
		if num > 0 {
			for _, l := range lists {
				resp, err := usdtServer(l.Address)
				if err != nil {
					ethSleep()
					continue
				}
				for _, tran := range resp {
					if strings.ToLower(tran.ReferenceAddress) == strings.ToLower(l.Address) {
						usdtDiposit(l.Uid, *tran)
					}
				}
				ethSleep()
			}
			offset += 100
		} else {
			offset = 0
			ethSleep()
		}
	}
}
