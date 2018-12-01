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
		Share: 1, //转让 不论是否已计算推广收益均不再计入
		Class: locked.Class,
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
