package job

import (
	"kuangchi_backend/models"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"strings"
	"time"
	"github.com/astaxie/beego/httplib"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"errors"
	"encoding/json"
	"math/big"
	"kuangchi_backend/utils"
	"math"
)

type etherscanEventLog struct {
	Status  string
	Message string
	Result  []json.RawMessage
}

func tokenServer(addressHash string) (*etherscanEventLog, error) {
	var resp etherscanEventLog
	tokenServer := httplib.Get("http://api.etherscan.io/api").SetTimeout(10*time.Second, 10*time.Second)
	tokenServer.Param("module", "logs")
	tokenServer.Param("action", "getLogs")
	tokenServer.Param("apikey", "YourApiKeyToken")
	tokenServer.Param("fromBlock", utils.FromBlock)
	tokenServer.Param("toBlock", "last")
	tokenServer.Param("topic0", "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
	tokenServer.Param("topic0_2_opr", "and")
	tokenServer.Param("topic2", addressHash)
	if err := tokenServer.ToJSON(&resp); err != nil {
		beego.Error(err)
		return nil, err
	}
	if resp.Status == "0" {
		if resp.Message != "No records found" {
			beego.Warn(resp.Message)
		}
		return nil, errors.New(resp.Message)
	}
	return &resp, nil
}

func handleToken(uid int, log types.Log, currencies []*models.KcCurrency, recentNumber uint64) {
	for _, currency := range currencies {
		if strings.ToLower(currency.Contract) == strings.ToLower(log.Address.String()) {
			tokenDiposit(uid, currency.Currency, currency.Decimals, log, recentNumber)
		}
	}
}

func tokenDiposit(uid int, currency string, decimals int, log types.Log, recentNumber uint64) {
	xInt := common.BytesToHash(log.Data).Big()
	x := new(big.Float).SetInt(xInt)
	y := big.NewFloat(math.Pow10(decimals))
	amount, _ := new(big.Float).Quo(x, y).Float64()
	from := common.HexToAddress(log.Topics[1].String()).String()
	to := common.HexToAddress(log.Topics[2].String()).String()
	o := orm.NewOrm()
	o.Begin()
	record := models.KcUserAddressRecord{Uid: uid, Currency: currency, Hash: log.TxHash.String()}
	created, _, err := o.ReadOrCreate(&record, "Uid", "Currency", "Hash")
	if err != nil {
		beego.Error(err)
		o.Rollback()
		return
	}
	if record.Status != 0 {
		o.Rollback()
		return
	}
	if created {
		record.From = from
		record.To = to
		record.Amount = amount
	}

	var confirmation uint64 = 0
	if recentNumber > log.BlockNumber {
		confirmation = recentNumber - log.BlockNumber
	}
	if confirmation > 12 {
		record.Status = 2
		record.CheckTime = time.Now()
		if err := models.AddAmount(o, uid, amount, currency, "deposit"); err != nil {
			beego.Error(err)
			o.Rollback()
			return
		}
	}
	//统计累计充值数量
	userAddress := models.KcUserAddress{Uid: uid, Currency: currency}
	if o.Read(&userAddress, "Uid", "Currency") == nil {
		userAddress.AllAmount += amount
		if _, err := o.Update(&userAddress, "AllAmount"); err != nil {
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
	o.Commit()
}

func tokenSleep() {
	//time.Sleep(time.Millisecond * 100)
	time.Sleep(time.Second * 10)
}

func unmarshalLog(raw []byte) (log *types.Log, err error) {
	if err = json.Unmarshal(raw, &log); err != nil {
		temp := strings.Replace(string(raw), `"logIndex":"0x"`, `"logIndex":"0x0"`, 1)
		temp = strings.Replace(temp, `"transactionIndex":"0x"`, `"transactionIndex":"0x0"`, 1)
		if err = json.Unmarshal([]byte(temp), &log); err != nil {
			return nil, err
		}
	}
	return log, nil
}

func TokenJob() {
	var lists []*models.KcUserAddress
	var offset = 0
	var recentNumber uint64
	o := orm.NewOrm()
	currencies := models.GetCurrencies("all")
	beego.Debug("token job start")
	for {
		recentNumber = getRecentNumber()
		num, err := o.QueryTable("kc_user_address").Filter("currency", "ETH").Limit(100, offset).All(&lists, "uid", "address")
		if err != nil {
			beego.Error(err)
			time.Sleep(time.Second * 10)
			continue
		}
		if num > 0 {
			for _, l := range lists {
				address := common.HexToHash(strings.ToLower(l.Address)).String()
				resp, err := tokenServer(address)
				if err != nil {
					tokenSleep()
					continue
				}

				for _, raw := range resp.Result {
					log, err := unmarshalLog(raw)
					if err != nil {
						beego.Error(err)
					} else {
						if len(log.Topics) != 3 {
							continue
						}
						if log.Topics[2].String() == address {
							handleToken(l.Uid, *log, currencies, recentNumber)
						}
					}
				}
				tokenSleep()
			}
			offset += 100
		} else {
			offset = 0
			tokenSleep()
			currencies = models.GetCurrencies("all")
		}
	}
}

func getRecentNumber() uint64 {
	var resp struct {
		Status    string
		Message   string
		Result    string
	}
	tokenServer := httplib.Get("https://api.etherscan.io/api").SetTimeout(10*time.Second, 10*time.Second)
	tokenServer.Param("module", "proxy")
	tokenServer.Param("action", "eth_blockNumber")
	tokenServer.Param("apikey", "YourApiKeyToken")
	if err := tokenServer.ToJSON(&resp); err != nil {
		beego.Error(err)
		return 0
	}
	if resp.Status == "0" {
		beego.Error(resp.Message)
		return 0
	}
	return toUnit64(resp.Result)
}

func toUnit64(a string) uint64 {
	if len(a) < 2 {
		return 0
	}
	raw := a[2:]
	var dec uint64
	for _, byte := range []byte(raw) {
		nib := decodeNibble(byte)
		if nib == badNibble {
			return 0
		}
		dec *= 16
		dec += uint64(nib)
	}
	return dec
}

const badNibble = ^uint64(0)

func decodeNibble(in byte) uint64 {
	switch {
	case in >= '0' && in <= '9':
		return uint64(in - '0')
	case in >= 'A' && in <= 'F':
		return uint64(in - 'A' + 10)
	case in >= 'a' && in <= 'f':
		return uint64(in - 'a' + 10)
	default:
		return badNibble
	}
}
