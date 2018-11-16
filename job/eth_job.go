package job

import (
	"errors"
	"kuangchi_backend/models"
	"kuangchi_backend/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"time"
)

type ethTransaction struct {
	Hash          string
	From          string
	To            string
	Value         string
	Confirmations string
	IsError       string
}

type etherscan struct {
	Status  string
	Message string
	Result  []ethTransaction
}

func ethServer(address string) (*etherscan, error) {
	resp := etherscan{}
	ethserver := httplib.Get("http://api.etherscan.io/api").SetTimeout(10*time.Second, 10*time.Second)
	ethserver.Param("module", "account")
	ethserver.Param("action", "txlist")
	ethserver.Param("apikey", "YourApiKeyToken")
	ethserver.Param("sort", "desc")
	ethserver.Param("address", address)
	if err := ethserver.ToJSON(&resp); err != nil {
		beego.Warn(ethserver.String())
		ethSleep()
		return nil, err
	}
	if resp.Status == "0" {
		if resp.Message != "No transactions found" {
			beego.Warn(resp.Message)
		}
		return nil, errors.New(resp.Message)
	}
	return &resp, nil
}

func ethDiposit(uid int, trans ethTransaction) {
	o := orm.NewOrm()
	o.Begin()
	record := models.KcUserAddressRecord{Uid: uid, Currency: "ETH", Hash: trans.Hash}
	amount, _ := strconv.ParseFloat(trans.Value, 10)
	created, _, err := o.ReadOrCreate(&record, "Uid", "Currency", "Hash")
	if err != nil {
		beego.Error(err)
		o.Rollback()
		return
	}
	confirmation, _ := strconv.ParseInt(trans.Confirmations, 10, 64)
	if created {
		if trans.IsError == "1" {
			record.Status = 1
			record.CheckTime = time.Now()
		} else if confirmation > 12 {
			record.Status = 2
			record.CheckTime = time.Now()
			if err := models.AddAmount(o, uid, amount/utils.ETHdecimal, "ETH", "deposit"); err != nil {
				beego.Error(err)
				o.Rollback()
				return
			}
		}
		record.Uid = uid
		record.From = trans.From
		record.To = trans.To
		record.Amount = amount / utils.ETHdecimal
		record.Confirmations = int(confirmation)
		if _, err := o.Update(&record); err != nil {
			beego.Error(err)
			o.Rollback()
			return
		}
	} else {
		if record.Status == 0 {
			if trans.IsError == "1" {
				record.Status = 1
				record.CheckTime = time.Now()
			} else if confirmation > 12 {
				record.Status = 2
				record.CheckTime = time.Now()
				if err := models.AddAmount(o, uid, amount/utils.ETHdecimal, "ETH", "deposit"); err != nil {
					beego.Error(err)
					o.Rollback()
					return
				}
			}
			record.Confirmations = int(confirmation)
			if _, err := o.Update(&record); err != nil {
				beego.Error(err)
				o.Rollback()
				return
			}
		}
	}
	o.Commit()
}

func ethSleep() {
	time.Sleep(time.Second * 5)
}

func ETHJob() {
	var lists []*models.KcUserAddress
	var offset int = 0
	o := orm.NewOrm()
	beego.Warn("eth cron start")
	for {
		num, err := o.QueryTable("kc_user_address").Filter("currency", "ETH").Limit(100, offset).All(&lists, "uid", "address")
		if err != nil {
			beego.Error(err)
			time.Sleep(time.Second * 10)
			continue
		}
		if num > 0 {
			for _, l := range lists {
				resp, err := ethServer(l.Address)
				if err != nil {
					ethSleep()
					continue
				}
				for _, trans := range resp.Result {
					if strings.ToLower(trans.To) == strings.ToLower(l.Address) {
						ethDiposit(l.Uid, trans)
					}
				}
				ethSleep()
			}
			offset += 100
		} else {
			offset = 0
		}
		ethSleep()
	}
}
