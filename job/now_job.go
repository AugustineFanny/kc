package job

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"time"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"strings"
	"strconv"
	"kuangchi_backend/utils"
)

type btccomAddress struct {
	AddrStr     string   `json:"addrStr"`
	Balance     float64  `json:"balance"`
	UnconfirmedBalance float64  `json:"unconfirmedBalance"`
}

func btcAmountServer(address string) {
	resp := btccomAddress{}
	url := fmt.Sprintf("https://chain.bitcoinworld.com/insight-api/addr/%s/?noTxList=1", address)
	//url := fmt.Sprintf("https://insight.bitpay.com/api/addr/%s/?noTxList=1", address)
	btcserver := httplib.Get(url).SetTimeout(10*time.Second, 10*time.Second)
	if err := btcserver.ToJSON(&resp); err != nil {
		beego.Warn(address)
		beego.Warn(btcserver.String())
		time.Sleep(time.Second)
		return
	}
	o := orm.NewOrm()
	params := orm.Params{"NowAmount": resp.Balance, "UnconfirmedAmount": resp.UnconfirmedBalance, "UpdateTime": time.Now()}
	o.QueryTable("kc_user_address").Filter("Address", resp.AddrStr).Update(params)
}

func UpdateBTCNowAmount() {
	var lists orm.ParamsList
	var offset int = 0
	o := orm.NewOrm()
	for {
		num, err := o.QueryTable("kc_user_address").Filter("Currency", "BTC").Limit(50, offset).ValuesFlat(&lists, "address")
		if err != nil {
			beego.Error(err)
			time.Sleep(time.Second * 10)
			continue
		}
		if num > 0 {
			for _, l := range lists {
				btcAmountServer(l.(string))
				// API频率限制
				//time.Sleep(time.Second)
				time.Sleep(time.Second * 10)
			}
			offset += 50
		} else {
			offset = 0
			time.Sleep(time.Hour * 6)
		}
	}
}

type ethAmount struct {
	Status  string
	Message string
	Result  []map[string]string
}

func ethAmountServer(address string) {
	resp := ethAmount{}
	ethserver := httplib.Get("http://api.etherscan.io/api").SetTimeout(10*time.Second, 10*time.Second)
	ethserver.Param("module", "account")
	ethserver.Param("action", "balancemulti")
	ethserver.Param("apikey", "YourApiKeyToken")
	ethserver.Param("tag", "latest")
	ethserver.Param("address", address)
	if err := ethserver.ToJSON(&resp); err != nil {
		beego.Error(err)
		return
	}
	if resp.Status == "0" {
		beego.Error(resp.Message)
		return
	}
	for _, arg := range resp.Result {
		balance, err := strconv.ParseFloat(arg["balance"], 64)
		if err != nil {
			beego.Error(arg)
			continue
		}
		o := orm.NewOrm()
		o.QueryTable("kc_user_address").Filter("Address", arg["account"]).Update(orm.Params{"NowAmount": balance/utils.ETHdecimal, "UpdateTime": time.Now()})
	}
}

func UpdateETHNowAmount() {
	var lists orm.ParamsList
	var offset int = 0
	o := orm.NewOrm()
	for {
		num, err := o.QueryTable("kc_user_address").Filter("currency", "ETH").Limit(20, offset).ValuesFlat(&lists, "address")
		if err != nil {
			beego.Error(err)
			time.Sleep(time.Second * 10)
			continue
		}
		if num > 0 {
			strArray := make([]string, len(lists))
			for i, l := range lists {
				address := strings.ToLower(l.(string))
				strArray[i] = address
			}
			ethAmountServer(strings.Join(strArray, ","))
			//time.Sleep(time.Millisecond * 200)
			time.Sleep(time.Second * 10)
			offset += 20
		} else {
			offset = 0
			time.Sleep(time.Hour * 6)
		}
	}
}
