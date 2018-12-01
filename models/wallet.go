package models

import (
	"kuangchi_backend/result"
	"kuangchi_backend/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"fmt"
	"time"
	"strconv"
	"strings"
	"math/rand"
	"math"
	"github.com/spf13/cast"
)

func AddAmount(o orm.Ormer, uid int, amount float64, currency, desc string, remark ...string) error {
	coin := KcWallet{Uid: uid, Currency: currency}
	_, id, err := o.ReadOrCreate(&coin, "Uid", "Currency")
	if err != nil {
		beego.Error(err)
		return err
	}
	cur := KcWallet{Id: int(id)}
	if err := o.ReadForUpdate(&cur); err != nil {
		beego.Error(err)
		return err
	}
	cur.Amount += amount
	if _, err := o.Update(&cur, "Amount"); err != nil {
		beego.Error(err)
		return err
	}
	fundChange := KcFundChange{Uid: uid, Currency: currency, Amount: amount, Direction: 0, Desc: desc}
	if len(remark) > 0 {
		fundChange.Remark = remark[0]
	}
	if _, err := o.Insert(&fundChange); err != nil {
		beego.Error(err)
		return err
	}
	return nil
}
//utils.Sub(sellerWallet.LockAmount, trade.Amount)
func ReduceAmount(o orm.Ormer, uid int, amount float64, lock bool, currency, desc string) error {
	wallet := KcWallet{Uid: uid, Currency: currency}
	if err := o.ReadForUpdate(&wallet, "Uid", "Currency"); err != nil {
		return result.ErrCode(100402)
	}
	if lock {
		wallet.LockAmount -= amount
		if wallet.LockAmount < 0 {
			return result.ErrMsg(fmt.Sprintf("Locked in assets not enough %f %s", amount, currency))
		}
	}
	wallet.Amount -= amount
	if wallet.UsableAmount() < 0 {
		return result.ErrMsg(fmt.Sprintf("Active in assets not enough %f %s", amount, currency))
	}
	if _, err := o.Update(&wallet, "Amount", "LockAmount"); err != nil {
		beego.Error(err)
		return err
	}
	fundChange := KcFundChange{Uid: uid, Currency: currency, Amount: amount, Direction: 1, Desc: desc}
	if _, err := o.Insert(&fundChange); err != nil {
		beego.Error(err)
		return result.ErrCode(100102)
	}
	return nil
}

func LockAmount(o orm.Ormer, uid int, currency string, amount float64) error {
	wallet := KcWallet{Uid: uid, Currency: currency}
	if err := o.ReadForUpdate(&wallet, "Uid", "Currency"); err != nil {
		return result.ErrCode(100402)
	}
	wallet.LockAmount += amount
	if wallet.Amount - wallet.LockAmount < 0 {
		return result.ErrMsg(fmt.Sprintf("Active in assets not enough %f %s", amount, currency))
	}
	if _, err := o.Update(&wallet, "LockAmount"); err != nil {
		beego.Error(err)
		return err
	}
	return nil
}

func SetAddress(u *User, currency string) string {
	var ap KcAddressPool

	o := orm.NewOrm()
	o.Begin()
	err := o.Raw("SELECT * FROM kc_address_pool WHERE flag = 0 and currency = ? ORDER BY id LIMIT 1", currency).QueryRow(&ap)
	if err != nil {
		beego.Error(currency, "地址用完了")
		return ""
	}
	ap.Flag = 1
	if _, err = o.Update(&ap, "Flag"); err != nil {
		o.Rollback()
		return ""
	}
	ua := KcUserAddress{Uid: u.Id, Currency: currency, Address: ap.Address, AddressIndex: ap.AddressIndex}
	if _, err = o.Insert(&ua); err != nil {
		o.Rollback()
		return ""
	}
	o.Commit()
	return ap.Address
}

func GetAddress(u *User, currency string) string {
	var address KcUserAddress

	o := orm.NewOrm()
	if cur := GetCurrency(currency); cur != nil {
		if cur.Token == 1 {
			currency = "ETH"
		}
	}
	query := o.QueryTable("kc_user_address").Filter("uid", u.Id).Filter("currency", currency)
	count, _ := query.Count()
	if count == 0 {
		return SetAddress(u, currency)
	}
	if count > 1 {
		beego.Error(u.Id, currency, "多个地址")
		return ""
	}
	query.One(&address)
	return address.Address
}

func GetFinance(u *User) ([]*KcWallet, error) {
	var coins []*KcWallet
	o := orm.NewOrm()

	_, err := o.QueryTable("kc_wallet").Filter("uid", u.Id).All(&coins)
	if err != nil {
		beego.Error(err)
		return nil, result.ErrCode(100102)
	}
	for _, wallet := range coins {
		wallet.Amount = utils.ShowFloat(wallet.Amount, 6)
		wallet.LockAmount = utils.ShowFloat(wallet.LockAmount, 6)
		if wallet.LockAmount < 0 {
			beego.Error(fmt.Sprintf("%d lock_amount < 0", wallet.Uid))
			wallet.LockAmount = 0
		}
	}
	if len(coins) == 0 {
		coins = append(coins, &KcWallet{Currency: "BTC", Amount: 0, LockAmount: 0})
	}
	return coins, nil
}

func GetWallet(u *User, currency string) []*KcWallet {
	wallet := KcWallet{Uid: u.Id, Currency: currency}
	o := orm.NewOrm()
	o.Read(&wallet, "Uid", "Currency")
	wallet.Amount = utils.ShowFloat(wallet.Amount, 6)
	wallet.LockAmount = utils.ShowFloat(wallet.LockAmount, 6)
	return []*KcWallet{&wallet}
}

func ApplyWithdraw(u *User, currency, address string, amount, fee float64, feeEth int) (err error) {
	if len(address) < 30 {
		return result.ErrMsg("invalid address")
	}
	o := orm.NewOrm()
	o.Begin()
	record := KcUserAddressRecord{
		Uid: u.Id,
		Currency: currency,
		To: address,
		FeeAmount: fee,
		Direction: 1,
		Status: 3,
	}
	if feeEth == 1 {
		//使用ETH做手续费
		record.Amount = amount
		record.FeeCurrency = "ETH"
		if err = LockAmount(o, u.Id, "ETH", fee); err != nil {
			o.Rollback()
			return result.ErrCode(100404)
		}
	} else {
		record.Amount = amount - fee
		record.FeeCurrency = currency
	}
	if err = LockAmount(o, u.Id, currency, amount); err != nil {
		o.Rollback()
		return err
	}

	if _, err = o.Insert(&record); err != nil {
		o.Rollback()
		beego.Error(err)
		return result.ErrCode(100102)
	}
	o.Commit()
	message := "您的提现申请已成功,请耐心等待审核结果"
	messageEn := "You've successfully submitted the withdraw application, please wait for the verification result."
	messageKo := "현금 인출 신청 성공하셨습니다. 심사 결과를 조금만 기다리세요."
	messageJp := "引き出しの申し込みを提出完了、審査結果を待ちしてください"
	SetMessage(u.Id, message, messageEn, messageKo, messageJp, "")
	return nil
}

func TransferInStation(u *User, currency, des string, amount float64) (err error) {
	if currency == "BTC" || currency == "ETH" {
		return result.ErrMsg(fmt.Sprintf("%s do not support the transfer", currency))
	}
	desUser := GetUserByIdent(des)
	if desUser == nil {
		return result.ErrCode(100407)
	}
	if u.Id == desUser.Id {
		return result.ErrMsg("can't transfer to yourself")
	}
	o := orm.NewOrm()
	o.Begin()

	if err := ReduceAmount(o, u.Id, amount, false, currency, "instation"); err != nil {
		o.Rollback()
		return err
	}

	if err := AddAmount(o, desUser.Id, amount, currency, "instation"); err != nil {
		o.Rollback()
		beego.Error(err)
		return result.ErrCode(100102)
	}
	record := KcUserAddressRecord{
		Uid: u.Id,
		Currency: currency,
		Amount: amount,
		From: u.Username,
		To: desUser.Username,
		FeeAmount: 0,
		Direction: 2,
		Status: 2,
	}
	if _, err = o.Insert(&record); err != nil {
		o.Rollback()
		beego.Error(err)
		return result.ErrCode(100102)
	}
	o.Commit()
	message := fmt.Sprintf("划转 %f %s 至 %s 账户", amount, currency, des)
	messageEn := fmt.Sprintf("Send %f %s to %s", amount, currency, des)
	SetMessage(u.Id, message, messageEn, messageEn, messageEn, "")
	message = fmt.Sprintf("获得 %f %s，由 %s 账户划转", amount, currency, u.Username)
	message = fmt.Sprintf("Get %f %s，From %s", amount, currency, u.Username)
	SetMessage(desUser.Id, message, messageEn, messageEn, messageEn, "")
	//ChangeAdShow(u.Id, currency)
	return nil
}

func GetRate(uid int, code string) float64 {
	currency := GetCurrency(code)
	if currency != nil {
		return currency.FeeTrade
	}
	return 0.005
}

func UsableBalance(u *User, currency string, rates ...float64) float64 {

	o := orm.NewOrm()
	wallet := KcWallet{Uid: u.Id, Currency: currency}
	if err := o.Read(&wallet, "Uid", "Currency"); err != nil {
		return 0
	}
	res := wallet.UsableAmount()
	if len(rates) > 0 {
		res /= 1 + rates[0]
	}

	return utils.ShowFloat(res, 6)
}

func handleLocked(o orm.Ormer, u *User, currency string, amount float64, class int) (err error) {
	locked := KcLocked{
		Uid: u.Id,
		Currency: currency,
		Amount: amount,
		TotalAmount: amount,
		StartDate: time.Now(),
		Class: class,
	}
	if u.InviterId == 0 {
		//无需计算推广奖励
		locked.Share = 1
	}
	wallet := KcWallet{Uid: u.Id, Currency: locked.Currency}
	if err := o.ReadForUpdate(&wallet, "Uid", "Currency"); err != nil {
		return result.ErrCode(100402)
	}
	wallet.LockAmount += locked.Amount
	wallet.MiningAmount += locked.Amount
	if wallet.Amount - wallet.LockAmount < 0 {
		return result.ErrMsg(fmt.Sprintf("Active in assets not enough %f %s", locked.Amount, locked.Currency))
	}
	if _, err := o.Update(&wallet, "LockAmount", "MiningAmount"); err != nil {
		beego.Error(err)
		return err
	}
	if _, err := o.Insert(&locked); err != nil {
		beego.Error(err)
		return result.ErrCode(100102)
	}
	return nil
}

func Locked(u *User, currency string, amount float64) (err error) {
	o := orm.NewOrm()
	o.Begin()
	if err := handleLocked(o, u, currency, amount, 2); err != nil {
		o.Rollback()
		return err
	}
	o.Commit()
	return nil
}

func GetLocked(u *User, currency string) []*KcLocked {
	var list []*KcLocked
	o := orm.NewOrm()
	_, err := o.QueryTable("kc_locked").Filter("uid", u.Id).Filter("currency", currency).OrderBy("-id").All(&list)
	if err != nil {
		return []*KcLocked{}
	}
	return list
}

func GetMining(u *User, currency string, pageNo int64) *utils.Page {
	o := orm.NewOrm()
	query := o.QueryTable("kc_fund_change").Filter("uid", u.Id).Filter("currency", currency).Filter("desc__in", "mining", "share", "competition")
	cnt, _ := query.Count()
	page := utils.SetPage(pageNo, cnt, 20)
	var minings []*KcFundChange
	num, err := query.Limit(page.PageSize, page.Offset).OrderBy("-id").All(&minings)
	if err != nil {
		beego.Error(err)
		return nil
	} else if num >= 0 {
		page.List = minings
		return page
	}
	return nil
}

func GetMiningStat(u *User, currency string) orm.Params {
	var list orm.ParamsList
	res := orm.Params{}
	o := orm.NewOrm()
	sql := "SELECT SUM(amount) balance FROM kc_fund_change " +
		   "WHERE uid = ? AND currency = ? AND `desc` in (\"mining\", \"share\", \"competition\") AND TO_DAYS(NOW()) - TO_DAYS(create_time) = 1"
	num, err := o.Raw(sql, u.Id, currency).ValuesFlat(&list, "balance")
	if err == nil && num > 0 {
		res["yesterday"] = list[0]
	} else {
		res["yesterday"] = nil
	}
	sql = "SELECT SUM(amount) balance FROM kc_fund_change " +
		  "WHERE uid = ? AND currency = ? AND `desc` in (\"mining\", \"share\", \"competition\") AND DATE_FORMAT(create_time, \"%y%m\") = DATE_FORMAT(CURDATE(), \"%y%m\")"
	num, err = o.Raw(sql, u.Id, currency).ValuesFlat(&list, "balance")
	if err == nil && num > 0 {
		res["month"] = list[0]
	} else {
		res["month"] = nil
	}
	return res
}

func GetNodes(currency string) []orm.Params {
	var maps []orm.Params
	o := orm.NewOrm()
	query := `SELECT username, mining_amount FROM kc_wallet
			  LEFT JOIN user ON currency = ? AND kc_wallet.uid = user.id
              WHERE mining_amount >= 1 ORDER BY mining_amount DESC LIMIT 70`
	_, err := o.Raw(query, currency).Values(&maps)
	if err != nil {
		return []orm.Params{}
	}
	return maps
}

func TransferLocked(u *User, id int, des string) (err error) {
	desUser := GetUserByIdent(des)
	if desUser == nil {
		return result.ErrCode(100407)
	}
	if u.Id == desUser.Id {
		return result.ErrMsg("can't transfer to yourself")
	}
	o := orm.NewOrm()
	locked := KcLocked{Id: id}
	if err := o.Read(&locked); err != nil {
		return result.ErrCode(100102)
	}
	if locked.Uid != u.Id {
		return result.ErrCode(100107)
	}
	if locked.Status != 0 {
		return result.ErrCode(100410)
	}
	cur := GetCurrency(locked.Currency)
	if cur == nil {
		return result.ErrCode(100102)
	}
	o.Begin()
	//来源减
	locked.Status = 2
	if _, err := o.Update(&locked, "status"); err != nil {
		o.Rollback()
		beego.Error(err)
		return result.ErrCode(100102)
	}
	wallet := KcWallet{Uid: u.Id, Currency: locked.Currency}
	if err := o.ReadForUpdate(&wallet, "Uid", "Currency"); err != nil {
		o.Rollback()
		beego.Error(err)
		return result.ErrCode(100402)
	}
	wallet.Amount -= locked.Amount
	wallet.LockAmount -= locked.Amount
	wallet.MiningAmount -= locked.Amount
	if _, err := o.Update(&wallet, "Amount", "LockAmount", "MiningAmount"); err != nil {
		o.Rollback()
		beego.Error(err)
		return err
	}
	fundChange := KcFundChange{Uid: u.Id, Currency: locked.Currency, Amount: locked.Amount, Direction: 1, Desc: "inlocked"}
	if _, err := o.Insert(&fundChange); err != nil {
		o.Rollback()
		beego.Error(err)
		return result.ErrCode(100102)
	}
	//目标加
	desLocked := KcLocked{
		Uid: desUser.Id,
		Currency: locked.Currency,
		Amount: locked.Amount,
		TotalAmount: locked.Amount, //转让后均重新计算
		StartDate: locked.StartDate,
		ExpireDate: time.Now().AddDate(0, 0, 60),
		Share: 1, //转让 不论是否已计算推广收益均不再计入
		Class: locked.Class,
	}
	if utils.Period == utils.UpComping || utils.Period == utils.Footstone {
		desLocked.ExpireDate = utils.FootstoneExpireDate
	}
	desLockedId, err := o.Insert(&desLocked)
	if err != nil {
		beego.Error(err)
		o.Rollback()
		return result.ErrCode(100102)
	}
	if err := AddAmount(o, desUser.Id, locked.Amount, locked.Currency, "inlocked"); err != nil {
		o.Rollback()
		return result.ErrCode(100102)
	}
	desWallet := KcWallet{Uid: desUser.Id, Currency: locked.Currency}
	if err := o.ReadForUpdate(&desWallet, "Uid", "Currency"); err != nil {
		o.Rollback()
		return result.ErrCode(100402)
	}
	desWallet.LockAmount += locked.Amount
	desWallet.MiningAmount += locked.Amount
	if _, err := o.Update(&desWallet, "LockAmount", "MiningAmount"); err != nil {
		o.Rollback()
		beego.Error(err)
		return err
	}
	//记录
	record := KcUserAddressRecord{
		Uid: u.Id,
		Currency: locked.Currency,
		Amount: locked.Amount,
		From: u.Username,
		To: desUser.Username,
		FeeAmount: 0,
		Direction: 3,
		Status: 2,
		Desc: strconv.Itoa(id),
		Remark: fmt.Sprintf("%d", desLockedId),
	}
	if _, err = o.Insert(&record); err != nil {
		o.Rollback()
		beego.Error(err)
		return result.ErrCode(100102)
	}
	o.Commit()
	return nil
}

func getFETSubscriptionNum() float64 {
	class := -1
	if utils.Period == utils.Footstone {
		class = 0
	} else if utils.Period == utils.Angel {
		class =1
	}
	params := orm.ParamsList{}
	o := orm.NewOrm()
	sql := `SELECT amounts FROM (SELECT currency, SUM(cur_amount) amounts FROM kc_subscription WHERE status = 2 AND class= ? GROUP BY currency) m WHERE currency = "FET"`
	_, err := o.Raw(sql, class).ValuesFlat(&params, "amounts")
	if err != nil {
		beego.Error(err)
		return -1
	}
	if len(params) != 1 {
		//没有锁仓
		return 0
	}
	amounts, err := strconv.ParseFloat(params[0].(string), 64)
	if err != nil {
		beego.Error(err)
		return -1
	}
	return amounts
}

func getFETSubscriptionHandle(amounts float64, base string) (float64, float64) {
	cur := GetCurrency(base)
	if cur == nil {
		return 0, 0
	}
	exchange, price := getFETSubscriptionPrice(amounts, cur.BasePrice)
	return utils.ShowFloat(exchange, 2), price
}

func getFETSubscriptionPrice(amounts float64, basePrice float64) (float64, float64) {
	if utils.Period == utils.Footstone {
		return getFETSubscriptionPriceFootstone(amounts, basePrice)
	} else if utils.Period == utils.Angel {
		return getFETSubscriptionPriceAngel(amounts, basePrice)
	}
	return 0, 0
}

func getFETSubscriptionPriceFootstone(amounts float64, basePrice float64) (float64, float64) {
	level := int(amounts) / 3000000
	switch level {
	case 0:  return basePrice / 0.81, 0.81
	case 1:  return basePrice / 0.82, 0.82
	case 2:  return basePrice / 0.83, 0.83
	case 3:  return basePrice / 0.84, 0.84
	case 4:  return basePrice / 0.85, 0.85
	case 5:  return basePrice / 0.86, 0.86
	case 6:  return basePrice / 0.87, 0.87
	case 7:  return basePrice / 0.88, 0.88
	case 8:  return basePrice / 0.89, 0.89
	case 9:  return basePrice / 0.90, 0.90
	case 10: return basePrice / 0.91, 0.91
	case 11: return basePrice / 0.92, 0.92
	case 12: return basePrice / 0.93, 0.93
	case 13: return basePrice / 0.94, 0.94
	case 14: return basePrice / 0.95, 0.95
	case 15: return basePrice / 0.96, 0.96
	case 16: return basePrice / 0.97, 0.97
	case 17: return basePrice / 0.98, 0.98
	case 18: return basePrice / 0.99, 0.99
	case 19: return basePrice / 1.00, 1.00
	//return: -1代表认购完成
	case 20: return -1, -1
	default: return 0, 0
	}
	return 0, 0
}

func getFETSubscriptionPriceAngel(amounts float64, basePrice float64) (float64, float64) {
	level := int(amounts) / 3000000
	switch level {
	case 0:  return basePrice / 1.01, 1.01
	case 1:  return basePrice / 1.02, 1.02
	case 2:  return basePrice / 1.03, 1.03
	case 3:  return basePrice / 1.04, 1.04
	case 4:  return basePrice / 1.05, 1.05
	case 5:  return basePrice / 1.06, 1.06
	case 6:  return basePrice / 1.07, 1.07
	case 7:  return basePrice / 1.08, 1.08
	case 8:  return basePrice / 1.09, 1.09
	case 9:  return basePrice / 1.10, 1.10
	case 10: return basePrice / 1.11, 1.11
	case 11: return basePrice / 1.12, 1.12
	case 12: return basePrice / 1.13, 1.13
	case 13: return basePrice / 1.14, 1.14
	case 14: return basePrice / 1.15, 1.15
	case 15: return basePrice / 1.16, 1.16
	case 16: return basePrice / 1.17, 1.17
	case 17: return basePrice / 1.18, 1.18
	case 18: return basePrice / 1.19, 1.19
	case 19: return basePrice / 1.20, 1.20
	case 20: return basePrice / 1.21, 1.21
	case 21: return basePrice / 1.22, 1.22
	case 22: return basePrice / 1.23, 1.23
	case 23: return basePrice / 1.24, 1.24
	case 24: return basePrice / 1.25, 1.25
	case 25: return basePrice / 1.26, 1.26
	case 26: return basePrice / 1.27, 1.27
	case 27: return basePrice / 1.28, 1.28
	case 28: return basePrice / 1.29, 1.29
	case 29: return basePrice / 1.30, 1.30
	case 30: return basePrice / 1.31, 1.31
	case 31: return basePrice / 1.32, 1.32
	case 32: return basePrice / 1.33, 1.33
	case 33: return basePrice / 1.34, 1.34
	case 34: return basePrice / 1.35, 1.35
	case 35: return basePrice / 1.36, 1.36
	case 36: return basePrice / 1.37, 1.37
	case 37: return basePrice / 1.38, 1.38
	case 38: return basePrice / 1.39, 1.39
	case 39: return basePrice / 1.40, 1.40
	case 40: return basePrice / 1.41, 1.41
	case 41: return basePrice / 1.42, 1.42
	case 42: return basePrice / 1.43, 1.43
	case 43: return basePrice / 1.44, 1.44
	case 44: return basePrice / 1.45, 1.45
	case 45: return basePrice / 1.46, 1.46
	case 46: return basePrice / 1.47, 1.47
	case 47: return basePrice / 1.48, 1.48
	case 48: return basePrice / 1.49, 1.49
	case 49: return basePrice / 1.50, 1.50
		//return: -1代表认购完成
	case 50: return -1, -1
	default: return 0, 0
	}
	return 0, 0
}

func getFETQuota(amounts float64) float64 {
	level := int(amounts) / 10000000
	switch level {
	case 0:    return 10000
	case 1, 2: return 30000
	default:   return 50000
	}
	return 0
}

func getFETSubscribed(uid int) float64 {
	params := orm.ParamsList{}
	o := orm.NewOrm()
	sql := `SELECT sum(cur_amount) amounts, uid FROM kc_subscription WHERE uid = ? AND currency = "FET" GROUP BY uid`
	_, err := o.Raw(sql, uid).ValuesFlat(&params, "amounts")
	if err != nil {
		beego.Error(err)
		return -1
	}
	if len(params) > 1 {
		beego.Error("有问题")
		return -1
	}
	if len(params) == 0 {
		return 0
	}
	amounts, err := strconv.ParseFloat(params[0].(string), 64)
	if err != nil {
		beego.Error(err)
		return -1
	}
	return amounts
}

func getFETDayQuota() float64 {
	start := utils.StartDate
	duration := time.Now().Sub(start).Hours()
	days := int64(duration) / 24 + 1
	return float64(days * 1000000)
}

func GetOrders(u *User, pageNo int64) *utils.Page {
	o := orm.NewOrm()
	query := o.QueryTable("kc_subscription").Filter("uid", u.Id).OrderBy("-id")
	cnt, _ := query.Count()
	page := utils.SetPage(pageNo, cnt, 10)
	var orders []*KcSubscription
	num, err := query.Limit(page.PageSize, page.Offset).All(&orders)
	if err != nil {
		beego.Error(err)
		return nil
	} else if num >= 0 {
		page.List = orders
		return page
	}
	return nil
}

func GetOrder(u *User, order string) (*KcSubscription, error) {
	o := orm.NewOrm()
	subscription := KcSubscription{Uid: u.Id, Order: order}
	if err := o.Read(&subscription, "uid", "order"); err != nil {
		return nil, result.ErrCode(100410)
	}
	return &subscription, nil
}

func GetHashrate(u *User) map[string]float64 {
	sql := `SELECT FLOOR(amount) amounts, user.id uid, user.inviter_id, user.parents FROM kc_locked
			LEFT JOIN user ON kc_locked.uid = user.id
			WHERE kc_locked.currency = ? AND kc_locked.share = 0 AND kc_locked.amount >= 1 AND FIND_IN_SET(?, parents)`
	var params []orm.Params
	res := map[string]float64{"private": 0, "competition": 0}
	o := orm.NewOrm()
	if _, err := o.Raw(sql, "FET", u.Id).Values(&params); err != nil {
		beego.Error(err)
		return res
	}
	root := strings.Split(u.DoAsParents(), ",")
	rootLen := len(root)

	ss := ShareStruct{}
	ss.Init()
	for _, m := range params {
		parents := strings.Split(m["parents"].(string), ",")
		inviterId, err := strconv.Atoi(m["inviter_id"].(string))
		if err != nil {
			beego.Error(err)
			continue
		}

		if inviterId == u.Id {
			ss.AppendExtension(m["amounts"].(string))
		}

		if len(parents) == rootLen {
			uid := m["uid"].(string)
			ss.AddCompetition(uid, m["amounts"].(string))
		} else if len(parents) > rootLen {
			uid := parents[rootLen]
			ss.AddCompetition(uid, m["amounts"].(string))
		} else {
			beego.Error(u.DoAsParents(), m["parents"].(string))
		}
	}

	res["private"] = ss.ExtensionRate()
	res["competition"] = ss.CompetitionRate()
	return res
}
