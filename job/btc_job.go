package job

import (
	"errors"
	"fmt"
	"kuangchi_backend/models"
	"kuangchi_backend/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/orm"
	"time"
)

type output struct {
	Addresses []string
	Value     float64
}

type list struct {
	Hash          string
	Outputs       []output
	Confirmations int
}

type data struct {
	List []list
}

type btccom struct {
	ErrNo  int    `json:"err_no"`
	ErrMsg string `json:"err_msg"`
	Data   data
}

type btcTransaction struct {
	Hash          string
	To            string
	Value         float64
	Confirmations int
}

func btcServer(address string) (*data, error) {
	resp := btccom{}
	url := fmt.Sprintf("https://chain.api.btc.com/v3/address/%s/tx", address)
	btcserver := httplib.Get(url).SetTimeout(10*time.Second, 10*time.Second)
	if err := btcserver.ToJSON(&resp); err != nil {
		beego.Warn(btcserver.String())
		btcSleep()
		return nil, err
	}

	if resp.ErrNo != 0 {
		if resp.ErrMsg != "Resource Not Found" {
			beego.Warn(resp.ErrMsg)
		}
		return nil, errors.New(resp.ErrMsg)
	}
	return &resp.Data, nil
}

func btcDiposit(uid int, trans btcTransaction) {
	o := orm.NewOrm()
	o.Begin()
	record := models.KcUserAddressRecord{Uid: uid, Currency: "BTC", Hash: trans.Hash}
	created, _, err := o.ReadOrCreate(&record, "Uid", "Currency", "Hash")
	if err != nil {
		beego.Error(err)
		o.Rollback()
		return
	}
	if created {
		if trans.Confirmations > 1 {
			record.Status = 2
			record.CheckTime = time.Now()
			if err := models.AddAmount(o, uid, trans.Value/utils.BTCDecimal, "BTC", "deposit"); err != nil {
				beego.Error(err)
				o.Rollback()
				return
			}
		}
		record.Uid = uid
		record.To = trans.To
		record.Amount = trans.Value / utils.BTCDecimal
		if _, err := o.Update(&record); err != nil {
			beego.Error(err)
			o.Rollback()
			return
		}
	} else {
		if record.Status == 0 {
			if trans.Confirmations > 1 {
				record.Status = 2
				record.CheckTime = time.Now()
				if err := models.AddAmount(o, uid, trans.Value/utils.BTCDecimal, "BTC", "deposit"); err != nil {
					beego.Error(err)
					o.Rollback()
					return
				}
			}
			if _, err := o.Update(&record); err != nil {
				beego.Error(err)
				o.Rollback()
				return
			}
		}
	}

	o.Commit()
}

func toTrans(address string, list list) (bool, *btcTransaction) {
	for _, output := range list.Outputs {
		for _, addr := range output.Addresses {
			if addr == address {
				return true, &btcTransaction{list.Hash, addr, output.Value, list.Confirmations}
			}
		}
	}
	return false, nil
}

func btcSleep() {
	time.Sleep(time.Second * 10)
}

func BTCJob() {
	var lists []*models.KcUserAddress
	var offset int = 0
	o := orm.NewOrm()
	beego.Warn("btc cron start")
	for {
		num, err := o.QueryTable("kc_user_address").Filter("currency", "BTC").Limit(100, offset).All(&lists, "uid", "address")
		if err != nil {
			beego.Error(err)
			time.Sleep(time.Second * 10)
			continue
		}
		if num > 0 {
			for _, l := range lists {
				resp, err := btcServer(l.Address)
				if err != nil {
					btcSleep()
					continue
				}
				for _, list := range resp.List {
					if valid, trans := toTrans(l.Address, list); valid {
						btcDiposit(l.Uid, *trans)
					}
				}
				// API频率限制
				btcSleep()
			}
			offset += 100
		} else {
			offset = 0
		}
		btcSleep()
	}
}
