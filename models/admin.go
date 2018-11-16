package models

import (
	"fmt"
	"time"
	"kuangchi_backend/utils"
	"kuangchi_backend/result"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"errors"
	"strings"
	"github.com/astaxie/beego/validation"
	"github.com/spf13/cast"
)

func CreateCurrency(currency KcCurrency) (*int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(&currency)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func CreateAdmin(u *Admin) (*int64, error) {
	valid := validation.Validation{}
	if v := valid.MinSize(u.Password, 6, "password"); !v.Ok {
		return nil, errors.New("密码太短")
	}

	u.Salt = utils.RandomString(10)
	u.EncodePasswd()

	o := orm.NewOrm()
	id, err := o.Insert(u)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func AdminSignIn(username, password string) (*Admin, error) {
	o := orm.NewOrm()
	admin := Admin{Username: username}
	if err := o.Read(&admin, "Username"); err != nil {
		return nil, errors.New("用户未注册")
	}
	if !admin.ValidatePassword(password) {
		return nil, errors.New("密码错误")
	}
	if admin.Status != 0 {
		return nil, errors.New("用户无法登录")
	}
	admin.LastTime = time.Now()
	if _, err := o.Update(&admin, "LastTime"); err != nil {
		beego.Error(err)
		return nil, err
	}
	return &admin, nil
}

func ChangeAdminPassword(u *Admin, password string, newPassword string) (err error) {
	valid := validation.Validation{}
	if !u.ValidatePassword(password) {
		return errors.New("密码错误")
	}
	if v := valid.MinSize(newPassword, 6, "newPassword"); !v.Ok {
		return errors.New("密码太短")
	}

	u.Salt = utils.RandomString(10)
	u.Password = newPassword
	u.EncodePasswd()

	o := orm.NewOrm()
	_, err = o.Update(u, "Password", "Salt")
	if err != nil {
		beego.Error(err)
		return errors.New("操作失败")
	}
	return nil
}

func formatSort(sort, order string, allows[]string) string {
	allow := false
	for _, x := range allows {
		if sort == x {
			allow = true
		}
	}
	if allow == false {
		return ""
	}
	if order == "descending" {
		return "-" + sort
	}
	return sort
}

func GetUsers(pageNo int64, sort, order string, ident string, group_id int) *utils.Page {
	o := orm.NewOrm()
	cnt, _ := o.QueryTable("user").Count()
	page := utils.SetPage(pageNo, cnt)
	var users []*User
	query := o.QueryTable("user")
	if ident != "" {
		if strings.Contains(ident, "@") {
			query = query.Filter("email", ident)
		} else if len(ident) == 11 && (ident[0] >= '0' && ident[0] <= '9'){
			query = query.Filter("mobile", ident)
		} else if ident[0] >= '0' && ident[0] <= '9' {
			query = query.Filter("id", ident)
		} else {
			query = query.Filter("username", ident)
		}
	}
	if group_id != -1 {
		query = query.Filter("group_id", group_id)
	}
	if sort != "" && order != ""{
		exprs := formatSort(sort, order, []string{"id", "last_time"})
		if exprs != "" {
			query = query.OrderBy(exprs)
		}
	}
	num, err := query.Limit(page.PageSize, page.Offset).All(&users)
	if err != nil {
		beego.Error(err)
		return nil
	} else if num >= 0 {
		for _, user := range users {
			if utils.RedisClient.IsExist("blocked:" + user.Username) {
				user.Status = 2
			}
		}
		page.List = users
		return page
	}
	return nil
}

func GetRealNames(pageNo, uid int64) *utils.Page {
	o := orm.NewOrm()
	query := o.QueryTable("real_name")
	if uid != -1 {
		query = query.Filter("uid", uid)
	} else {
		query = query.Filter("status", 1)
	}
	cnt, _ := query.Count()
	page := utils.SetPage(pageNo, cnt)
	var realNames []*RealName
	num, err := query.Limit(page.PageSize, page.Offset).All(&realNames)
	if err != nil {
		beego.Error(err)
		return nil
	} else if num >= 0 {
		page.List = realNames
		return page
	}
	return nil
}

func ChangeRole(status StatusDto) (err error) {
	realName := RealName{Id: status.Id}
	o := orm.NewOrm()
	o.Begin()
	if err = o.Read(&realName); err != nil {
		o.Rollback()
		return err
	}
	realName.AuthTime = time.Now()
	realName.Status = status.Status
	realName.Desc = status.Desc
	if _, err = o.Update(&realName, "AuthTime", "Status", "Desc"); err != nil {
		o.Rollback()
		return err
	}
	user := User{Id: realName.Uid}
	if err = o.Read(&user); err != nil {
		o.Rollback()
		return err
	}
	if user.Role != 1 && user.Role != 2 {
		o.Rollback()
		return errors.New("该用户暂无法认证")
	}
	user.Role = status.Status
	if _, err = o.Update(&user, "Role"); err != nil {
		o.Rollback()
		return err
	}
	var message, messageEn , messageKo, messageJp string
	if status.Status == 2 {
		message = fmt.Sprintf("身份验证审核未通过，原因: %s", status.Desc)
		messageEn = fmt.Sprintf("Failed the ID authentication, specific reasons are: %s", status.Desc)
		messageKo = fmt.Sprintf("신분 인증 심사에 실패하셨습니다. 그 원인은 아래와 같습니다: %s", status.Desc)
		messageJp = fmt.Sprintf("個人情報審査失敗、原因：%s", status.Desc)
	} else {
		message = fmt.Sprintf("身份验证审核通过")
		messageEn = fmt.Sprintf("Passed the ID authentication.")
		messageKo = fmt.Sprintf("신분 인증 심사에 통과하셨습니다.")
		messageJp = fmt.Sprintf("個人情報審査完了")
	}
	o.Commit()
	SetMessage(realName.Uid, message, messageEn, messageKo, messageJp, "")
	return nil
}

func GetKycs(pageNo, uid int64) *utils.Page {
	o := orm.NewOrm()
	query := o.QueryTable("kyc")
	if uid != -1 {
		query = query.Filter("uid", uid)
	} else {
		query = query.Filter("status", 1)
	}
	cnt, _ := query.Count()
	page := utils.SetPage(pageNo, cnt)
	var kycs []*Kyc
	num, err := query.Limit(page.PageSize, page.Offset).All(&kycs)
	if err != nil {
		beego.Error(err)
		return nil
	} else if num >= 0 {
		page.List = kycs
		return page
	}
	return nil
}

func ChangeKyc(status StatusDto) (err error) {
	kyc := Kyc{Id: status.Id}
	o := orm.NewOrm()
	o.Begin()
	if err = o.Read(&kyc); err != nil {
		o.Rollback()
		return err
	}
	kyc.AuthTime = time.Now()
	kyc.Status = status.Status
	kyc.Desc = status.Desc
	if _, err = o.Update(&kyc, "AuthTime", "Status", "Desc"); err != nil {
		o.Rollback()
		return err
	}
	user := User{Id: kyc.Uid}
	if err = o.Read(&user); err != nil {
		o.Rollback()
		return err
	}
	if user.Kyc != 1 && user.Kyc != 2 {
		o.Rollback()
		return errors.New("该用户暂无法认证")
	}
	user.Kyc = status.Status
	if _, err = o.Update(&user, "Kyc"); err != nil {
		o.Rollback()
		return err
	}
	var message string
	if status.Status == 2 {
		message = fmt.Sprintf("KYC审核未通过，原因: %s", status.Desc)
	} else {
		message = fmt.Sprintf("KYC审核通过")
	}
	o.Commit()
	SetMessage(kyc.Uid, message, "", "", "", "")
	return nil
}

func GetTransfers(uid int, currency string, status, direction int, pageNo int64, pageSize ...int64) *utils.Page {
	o := orm.NewOrm()
	query := o.QueryTable("kc_user_address_record")
	if uid != -1 {
		query = query.Filter("uid", uid)
	}
	if currency != "" {
		query = query.Filter("currency", currency)
	}
	if status != -1 {
		query = query.Filter("status", status)
	}
	if direction != -1 {
		query = query.Filter("direction", direction)
	}
	cnt, _ := query.Count()
	page := utils.SetPage(pageNo, cnt, pageSize...)
	var records []*KcUserAddressRecord
	num, err := query.OrderBy("-id").Limit(page.PageSize, page.Offset).All(&records)
	if err != nil {
		beego.Error(err)
		return nil
	} else if num >= 0 {
		page.List = records
		return page
	}
	return nil
}

func ThawAmount(o orm.Ormer, uid int, currency string, amount float64) error {
	wallet := KcWallet{Uid: uid, Currency: currency}
	if err := o.ReadForUpdate(&wallet, "Uid", "Currency"); err != nil {
		return result.ErrCode(100402)
	}
	wallet.LockAmount -= amount
	if wallet.LockAmount < 0 {
		return result.ErrMsg("锁币金额小于0")
	}
	if _, err := o.Update(&wallet, "LockAmount"); err != nil {
		beego.Error(err)
		return err
	}
	return nil
}

func DeductAmount(o orm.Ormer, uid int, currency string, amount float64) error {
	wallet := KcWallet{Uid: uid, Currency: currency}
	if err := o.ReadForUpdate(&wallet, "Uid", "Currency"); err != nil {
		return result.ErrCode(100402)
	}
	wallet.LockAmount -= amount
	if wallet.LockAmount < 0 {
		return result.ErrMsg("锁币金额小于0")
	}
	wallet.Amount -= amount
	if wallet.Amount < 0 {
		return result.ErrMsg("金额小于0")
	}
	if _, err := o.Update(&wallet, "Amount", "LockAmount"); err != nil {
		beego.Error(err)
		return err
	}
	fundChange := KcFundChange{Uid: uid, Currency: currency, Amount: amount, Direction: 1, Desc: "withdraw"}
	if _, err := o.Insert(&fundChange); err != nil {
		beego.Error(err)
		return err
	}
	return nil
}

func TransferAuth(status StatusDto) (err error) {
	transfer := KcUserAddressRecord{Id: status.Id}
	o := orm.NewOrm()
	o.Begin()
	if err = o.Read(&transfer); err != nil {
		o.Rollback()
		return err
	}
	if transfer.Status != 3 {
		o.Rollback()
		return errors.New("非待审核状态")
	}
	if transfer.Direction != 1 {
		o.Rollback()
		return errors.New("非提币条目")
	}
	user := User{Id: transfer.Uid}
	if err = o.Read(&user); err != nil {
		o.Rollback()
		return err
	}
	transfer.CheckTime = time.Now()
	transfer.Status = status.Status
	if status.Status == 2 {
		transfer.Hash = status.Desc
		if _, err = o.Update(&transfer, "CheckTime", "Status", "Hash"); err != nil {
			o.Rollback()
			return err
		}
	} else {
		transfer.Remark = status.Desc
		if _, err = o.Update(&transfer, "CheckTime", "Status", "Remark"); err != nil {
			o.Rollback()
			return err
		}
	}

	var message, messageEn, messageKo, messageJp string
	if status.Status == 4 {
		if err := ThawAmount(o, user.Id, transfer.Currency, transfer.AllAmount()); err != nil {
			o.Rollback()
			return err
		}
		if transfer.Currency != transfer.FeeCurrency {
			if err := ThawAmount(o, user.Id, transfer.FeeCurrency, transfer.FeeAmount); err != nil {
				o.Rollback()
				return err
			}
		}
		message = fmt.Sprintf("您的账户提现%f %s, 审核未通过，原因: %s", transfer.Amount, transfer.Currency, status.Desc)
		messageEn = fmt.Sprintf("Application failed to withdraw %f %s, reasons are: %s", transfer.Amount, transfer.Currency, status.Desc)
		messageKo = fmt.Sprintf("당신의 계정에서 %f %s 현금 인출하려는데 실패하셨습니다. 그 원인은 아래와 같습니다: %s", transfer.Amount, transfer.Currency, status.Desc)
		messageJp = fmt.Sprintf("%f %sの引き出しの審査失敗、原因： %s", transfer.Amount, transfer.Currency, status.Desc)
	}

	if status.Status == 2 {
		if err := DeductAmount(o, user.Id, transfer.Currency, transfer.AllAmount()); err != nil {
			o.Rollback()
			return err
		}
		if transfer.Currency != transfer.FeeCurrency {
			if err := DeductAmount(o, user.Id, transfer.FeeCurrency, transfer.FeeAmount); err != nil {
				o.Rollback()
				return err
			}
		}
		message = fmt.Sprintf("您的账户提现%f %s, 已提现成功", transfer.Amount, transfer.Currency)
		messageEn = fmt.Sprintf("Successfully withdrawn %f %s from your account.", transfer.Amount, transfer.Currency)
		messageKo = fmt.Sprintf("당신의 계정에서 %f %s 현금 인출 성공하셨습니다.", transfer.Amount, transfer.Currency)
		messageJp = fmt.Sprintf("%f %sの引き出し成功", transfer.Amount, transfer.Currency)
	}
	o.Commit()
	SetMessage(transfer.Uid, message, messageEn, messageKo, messageJp, "")
	return nil
}

func SetMessage(uid int, message, messageEn, messageKo, messageJp, extra string, messageType ...string) {
	o := orm.NewOrm()
	msg := KcMessage{Uid: uid, Message: message, MessageEn: messageEn, MessageKo: messageKo, MessageJp: messageJp, Extra: extra, MessageType: "default"}
	if len(messageType) > 0 {
		msg.MessageType = messageType[0]
	}
	if _, err := o.Insert(&msg); err != nil {
		//message插入失败不回滚
		beego.Error(err)
	}
}

func LogOperation(u *Admin,
	api string,
	args interface{},
	sheet string,
	rowId int,
	status bool,
	remark ...interface{}) {
	log := AdminOperationLog{Aid:u.Id, UserName: u.Username, Api: api, Sheet: sheet, RowId: rowId, Status: status}
	log.Args = fmt.Sprintf("%#v", args)
	if len(remark) > 0 {
		log.Remark = fmt.Sprintf("%#v", remark)
	}

	o := orm.NewOrm()
	if _, err := o.Insert(&log); err != nil {
		beego.Error(err)
	}
}

func GetAdmins() []*Admin {
	var list []*Admin
	o := orm.NewOrm()

	_, err := o.QueryTable("admin").All(&list)
	if err != nil {
		beego.Error(err)
		return nil
	}
	return list
}

func SetAdmin(admin Admin) error {
	o := orm.NewOrm()
	if _, err := o.Update(&admin, "Username", "Status"); err != nil {
		return err
	}
	return nil
}

func SetAdminMobile(admin Admin) error {
	o := orm.NewOrm()
	if _, err := o.Update(&admin, "Mobile"); err != nil {
		return err
	}
	return nil
}

func SetCurrency(currency KcCurrency) error {
	o := orm.NewOrm()
	if _, err := o.Update(&currency, "MinLock", "Token", "Contract", "ConfirmNum", "Recharge", "Withdraw",
									 "TransFlag", "Decimals", "BasePrice", "MInterestRate", "CInterestRate"); err != nil {
		return err
	}
	return nil
}

func GetLogs(pageNo int64, ident, api string, status int) *utils.Page {
	o := orm.NewOrm()
	query := o.QueryTable("admin_operation_log")
	if ident != "" {
		query = query.Filter("username", ident)
	}
	if api != "" {
		query = query.Filter("api", api)
	}
	if status != -1 {
		query = query.Filter("status", status)
	}
	cnt, _ := query.Count()
	page := utils.SetPage(pageNo, cnt)
	var logs []*AdminOperationLog
	num, err := query.OrderBy("-id",).Limit(page.PageSize, page.Offset).All(&logs)
	if err != nil {
		beego.Error(err)
		return nil
	} else if num >= 0 {
		page.List = logs
		return page
	}
	return nil
}

func BackGetWallets(pageNo int64, ident, currency string, noZero bool, sort, order string) *utils.Page {
	o := orm.NewOrm()
	query := o.QueryTable("kc_wallet")
	if ident != "" {
		query = query.Filter("uid", ident)
	}
	if currency != "" {
		query = query.Filter("currency", currency)
	}
	query = query.Filter("Amount__gt", 0)
	if noZero {
		query = query.Filter("MiningAmount__gt", 0)
	}
	if sort != "" {
		if order == "descending" {
			sort = "-" + sort
		}
		query = query.OrderBy(sort)
	}
	cnt, _ := query.Count()
	page := utils.SetPage(pageNo, cnt)
	var wallets []*KcWallet
	num, err := query.Limit(page.PageSize, page.Offset).All(&wallets)
	if err != nil {
		beego.Error(err)
		return nil
	} else if num >= 0 {
		page.List = wallets
		return page
	}
	return nil
}

func BackGetAddresses(pageNo int64, uid int, currency string, noZero bool) *utils.Page {
	o := orm.NewOrm()
	query := o.QueryTable("kc_user_address")
	if uid != -1 {
		query = query.Filter("uid", uid)
	}
	if currency != "" {
		query = query.Filter("currency", currency)
	}
	if noZero == true {
		query = query.Filter("NowAmount__gt", 0)
	}
	cnt, _ := query.Count()
	page := utils.SetPage(pageNo, cnt)
	var addresses []*KcUserAddress
	num, err := query.OrderBy("-id").Limit(page.PageSize, page.Offset).All(&addresses)
	if err != nil {
		beego.Error(err)
		return nil
	} else if num >= 0 {
		page.List = addresses
		return page
	}
	return nil
}

func BackGetAddressePool(pageNo int64, currency string, flag int) *utils.Page {
	o := orm.NewOrm()
	query := o.QueryTable("kc_address_pool")
	if currency != "" {
		query = query.Filter("currency", currency)
	}
	if flag != -1 {
		query = query.Filter("flag", flag)
	}
	cnt, _ := query.Count()
	page := utils.SetPage(pageNo, cnt)
	var pool []*KcAddressPool
	num, err := query.OrderBy("id").Limit(page.PageSize, page.Offset).All(&pool)
	if err != nil {
		beego.Error(err)
		return nil
	} else if num >= 0 {
		page.List = pool
		return page
	}
	return nil
}

func Stat(method string) []orm.Params {
	var maps []orm.Params
	o := orm.NewOrm()
	var sql string
	switch method {
	case "balance":
		sql = "SELECT currency, SUM(amount) balance FROM kc_wallet GROUP BY currency"
	case "recharge":
		sql = "SELECT currency, SUM(amount) balance FROM kc_user_address_record WHERE direction = 0 GROUP BY currency"
	case "withdraw":
		sql = "SELECT currency, SUM(amount) balance FROM kc_user_address_record WHERE direction = 1 AND status = 2 GROUP BY currency"
	case "sign_up_today":
		sql = "SELECT count(1) num FROM user WHERE TO_DAYS(create_time) = TO_DAYS(NOW())"
	case "sign_up_yesterday":
		sql = "SELECT count(1) num FROM user WHERE TO_DAYS(create_time) = TO_DAYS(NOW()) - 1"
	case "sign_up_month":
		sql = `SELECT count(1) num FROM user WHERE DATE_FORMAT(create_time, "%y%m") = DATE_FORMAT(NOW(), "%y%m")`
	case "sign_up_prev_month":
		sql = `SELECT count(1) num FROM user WHERE PERIOD_DIFF(DATE_FORMAT(NOW(), "%y%m"), DATE_FORMAT(create_time, "%y%m")) = 1`
	case "fet_stat":
		sql = `SELECT SUM(amount) amounts, SUM(mining_amount) mining_amounts FROM kc_wallet WHERE currency = "FET" GROUP BY currency`
	default:
		return []orm.Params{}
	}
	num, err := o.Raw(sql).Values(&maps)
	if err == nil && num > 0 {
		return maps
	}
	return []orm.Params{}
}

func BackGetInvites(startDate, endDate string) []orm.Params {
	var maps []orm.Params
	o := orm.NewOrm()
	if startDate == "" {
		startDate = "2018-01-01 00:00:00"
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02 15:04:05")
	}
	sql := `
		SELECT y.*, user.username FROM (
			SELECT uid,
				MAX(CASE grade WHEN 0 THEN num ELSE 0 END) m1,
				MAX(CASE grade WHEN 1 THEN num ELSE 0 END) m2,
				MAX(CASE grade WHEN 2 THEN num ELSE 0 END) m3,
				MAX(CASE grade WHEN 3 THEN num ELSE 0 END) m4,
				MAX(CASE grade WHEN 4 THEN num ELSE 0 END) m5,
				MAX(CASE grade WHEN 5 THEN num ELSE 0 END) m6
				FROM (SELECT count(1) num, uid, grade FROM invite WHERE create_time > ? AND create_time <= ? GROUP BY uid, grade) x GROUP BY uid
			) y
		LEFT JOIN user ON y.uid = user.id
	`
	num, err := o.Raw(sql, startDate, endDate).Values(&maps)
	if err == nil && num > 0 {
		return maps
	}
	return []orm.Params{}
}

func BackGetChildAmounts() []orm.Params {
	var maps []orm.Params
	o := orm.NewOrm()
	sql := `
		SELECT username,
			   invite.uid,
			   SUM(amount) amount,
			   SUM(mining_amount) lock_amount,
			   SUM(if(grade = 0, amount, 0)) m1_amount,
			   SUM(if(grade = 0, mining_amount, 0)) m1_lock_amount
		FROM invite
		INNER JOIN kc_wallet ON invite.invite = kc_wallet.uid
		LEFT JOIN user ON invite.uid = user.id
		WHERE currency = "FET"
		GROUP BY invite.uid
	`
	num, err := o.Raw(sql).Values(&maps)
	if err == nil && num > 0 {
		return maps
	}
	return []orm.Params{}
}

func GetGroups() []*Group {
	o := orm.NewOrm()
	var groups []*Group
	_, err := o.QueryTable("group").All(&groups)
	if err != nil {
		beego.Error(err)
		return nil
	}
	return groups
}

func ChangeUser(user User) error {
	o := orm.NewOrm()
	if _, err := o.Update(&user, "GroupId"); err != nil {
		return err
	}
	return nil
}

func CreateGroup(name string) (*int64, error) {
	o := orm.NewOrm()
	group := Group{Name: name}
	id, err := o.Insert(&group)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func BackSendCaptcha(mobile, captchaType string) error {
	if _, err := utils.RedisClient.GetString("HOLD:" + captchaType + ":" + mobile); err == nil {
		return result.ErrCode(100112)
	}

	captcha := utils.RandomCaptcha(5)

	if err := utils.RedisClient.Set("captcha:" + captchaType + ":" + mobile, captcha, 24 * 60 * 60); err != nil {
		return result.ErrCode(100304)
	}

	utils.RedisClient.Set("HOLD:" + captchaType + ":" + mobile, "hold", 50)

	go utils.SendSMS(mobile, captcha)
	return nil
}

func SetSuperPassword(u *Admin, newPassword string) (err error) {
	valid := validation.Validation{}
	if v := valid.MinSize(newPassword, 6, "newPassword"); !v.Ok {
		return errors.New("密码太短")
	}

	u.FundSalt = utils.RandomString(10)
	u.FundPassword = newPassword
	u.EncodeFundPasswd()

	o := orm.NewOrm()
	_, err = o.Update(u, "FundPassword", "FundSalt")
	if err != nil {
		beego.Error(err)
		return errors.New("操作失败")
	}
	return nil
}

func UserRecharge(uid int, currency string ,amount float64) error {
	o := orm.NewOrm()
	o.Begin()
	if err := AddAmount(o, uid, amount, currency, "distribute"); err != nil {
		o.Rollback()
		return err
	}
	//如果是FET，直接锁仓
	if currency == "FET" {
		user := GetUserById(uid)
		switch utils.Period {
		case utils.UpComping, utils.Footstone:
			expire := utils.FootstoneExpireDate
			if err := handleLocked(o, user, currency, amount, 0, expire); err != nil {
				o.Rollback()
				return err
			}
		case utils.Angel:
			expire := time.Now().AddDate(0, 0, 60)
			if err := handleLocked(o, user, currency, amount, 0, expire); err != nil {
				o.Rollback()
				return err
			}
		}
	}
	o.Commit()
	message := fmt.Sprintf("获得 %f %s", amount, currency)
	messageEn := fmt.Sprintf("Get %f %s", amount, currency)
	messageKo := fmt.Sprintf("얻다 %f %s", amount, currency)
	messageJp := fmt.Sprintf("獲得 %f %s", amount, currency)
	SetMessage(uid, message, messageEn, messageKo, messageJp, "")
	return nil
}

func BackGetFundChanges(pageNo int64, uid int, currency string, direction int, desc string) *utils.Page {
	o := orm.NewOrm()
	query := o.QueryTable("kc_fund_change")
	if uid != -1 {
		query = query.Filter("uid", uid)
	}
	if currency != "" {
		query = query.Filter("currency", currency)
	}
	if direction != -1 {
		query = query.Filter("direction", direction)
	}
	if desc != "" {
		query = query.Filter("desc", desc)
	}
	cnt, _ := query.Count()
	page := utils.SetPage(pageNo, cnt)
	var funds []*KcFundChange
	num, err := query.OrderBy("-id").Limit(page.PageSize, page.Offset).All(&funds)
	if err != nil {
		beego.Error(err)
		return nil
	} else if num >= 0 {
		page.List = funds
		return page
	}
	return nil
}

func BatchSms(content, mobile string) (err error) {
	if content == "" {
		return result.ErrMsg("参数不能为空 content")
	}
	if mobile != "" {
		res := utils.PushSms(content + "【投币时代】", mobile)
		return res
	}
	var list orm.ParamsList
	o := orm.NewOrm()
	_, err =  o.QueryTable("user").ValuesFlat(&list, "mobile")
	if err != nil {
		return err
	}
	for _, mobile := range list {
		des := mobile.(string)
		if des != "" {
			utils.PushSms(content + "【投币时代】", des)
		}
	}
	return nil
}

func BatchEmail(content, email string) (err error) {
	if content == "" {
		return result.ErrMsg("参数不能为空 content")
	}
	if email != "" {
		go utils.SendContentEmail(email, "FADAX", content)
		return nil
	}
	var list orm.ParamsList
	o := orm.NewOrm()
	_, err =  o.QueryTable("user").ValuesFlat(&list, "email")
	if err != nil {
		return err
	}
	for _, email := range list {
		des := email.(string)
		if des != "" {
			go utils.SendContentEmail(des, "FADAX", content)
		}
	}
	return nil
}

func SendCaptchaToUserByUsername(username string) error {
	user := GetUserByUsername(username)
	if user == nil {
		return errors.New("查无此用户")
	}
	if user.Verified("mobile") {
		return SendSMS(user.Mobile, "MANAGE", user.CountryCode)
	} else if user.Verified("email") {
		return SendEmail(user.Email, "MANAGE")
	} else {
		return errors.New("无手机号无邮箱？？？")
	}
	return nil
}

func SendSMSToAdmin(aid int, content string) {
	o := orm.NewOrm()
	admin := Admin{Id: aid}
	if err := o.Read(&admin); err != nil {
		return
	}
	if admin.Mobile != "" {
		go utils.PushSms(content + "【投币时代】", admin.Mobile)
	}
}

func SendSMSToMultiAdmin(content string, aids ...int) {
	if len(aids) != 0 {
		for _, aid := range aids {
			SendSMSToAdmin(aid, content)
		}
	} else {
		o := orm.NewOrm()
		var list orm.ParamsList
		o.QueryTable("admin").ValuesFlat(&list, "mobile")
		for _, mobile := range list {
			if mobile != "" {
				go utils.PushSms(content + "【投币时代】", mobile.(string))
			}
		}
	}
}

const (
	Instations = 2
	Inlocked = 3
)

func GetInstations(username, currency string, direction int, in int, pageNo int64, pageSize ...int64) *utils.Page {
	o := orm.NewOrm()
	query := o.QueryTable("kc_user_address_record")
	if username != "" {
		if direction == 0 {
			query = query.Filter("from", username)
		} else if direction == 1 {
			query = query.Filter("to", username)
		} else {
			//todo
			query = query.Filter("from", username)
		}
	}
	if currency != "" {
		query = query.Filter("currency", currency)
	}
	query = query.Filter("direction", in)
	cnt, _ := query.Count()
	page := utils.SetPage(pageNo, cnt, pageSize...)
	var records []*KcUserAddressRecord
	num, err := query.OrderBy("-id").Limit(page.PageSize, page.Offset).All(&records)
	if err != nil {
		beego.Error(err)
		return nil
	} else if num >= 0 {
		page.List = records
		return page
	}
	return nil
}

func GetProfitForDate(desc, date string) (maps []orm.Params) {
	o := orm.NewOrm()
	sql := `SELECT km.*, user.username FROM kc_mining km LEFT JOIN user ON km.uid = user.id
			WHERE description = ? AND DATE_FORMAT(km.create_time, "%Y-%m-%d") = ? ORDER BY reward DESC`
	num, err := o.Raw(sql, desc, date).Values(&maps)
	if err == nil && num > 0{
		return maps
	}
	return []orm.Params{}
}

func GetProfitForMonth(desc, startDate, endDate string) (maps []orm.Params) {
	o := orm.NewOrm()
	sql := "SELECT m.*, user.username FROM " +
				"(SELECT kfc.uid, kfc.currency, SUM(kfc.amount) amount FROM kc_fund_change kfc " +
				"WHERE `desc` = ? AND kfc.create_time > ? AND kfc.create_time <= ? " +
				"GROUP BY kfc.uid, kfc.currency ORDER BY amount DESC) m " +
			"LEFT JOIN user ON m.uid = user.id"
	num, err := o.Raw(sql, desc, startDate, endDate).Values(&maps)
	if err == nil && num > 0{
		return maps
	}
	return []orm.Params{}
}

func BackGetSubscriptions(pageNo int64, ident, order, base string, status int, date string) *utils.Page {
	o := orm.NewOrm()
	query := o.QueryTable("kc_subscription")
	if ident != "" {
		uid := 0
		if strings.Contains(ident, "@") {
			user := GetUserByEmail(ident)
			if user != nil {
				uid = user.Id
			}
		} else if ident[0] < '0' || ident[0] > '9' {
			user := GetUserByUsername(ident)
			if user != nil {
				uid = user.Id
			}
		} else {
			uid = cast.ToInt(ident)
		}
		query = query.Filter("uid", uid)
	}
	if order != "" {
		query = query.Filter("order", order)
	}
	if base != "" {
		query = query.Filter("base", base)
	}
	if status != -1 {
		query = query.Filter("status", status)
	}
	if date != "" {
		if date, err := time.Parse("2006-01-02", date); err == nil {
			query = query.Filter("create_time__gte", date).Filter("create_time__lt", date.AddDate(0, 0,1))
		}
	}
	cnt, _ := query.Count()
	page := utils.SetPage(pageNo, cnt, 50)
	var kcSubscription []*KcSubscription
	num, err := query.OrderBy("-id").Limit(page.PageSize, page.Offset).All(&kcSubscription)
	if err != nil {
		beego.Error(err)
		return nil
	} else if num >= 0 {
		page.List = kcSubscription
		return page
	}
	return nil
}

func BackGetSubmissions(order string) ([]*KcSubscriptionSubmission, error) {
	var submissions []*KcSubscriptionSubmission
	o := orm.NewOrm()
	_, err := o.QueryTable("kc_subscription_submission").Filter("order", order).OrderBy("-id").All(&submissions)
	if err != nil {
		beego.Error(err)
		return []*KcSubscriptionSubmission{}, err
	}
	return submissions, nil
}

func ConfirmOrder(id, status int, remark string) error {
	if status == 2 && remark != "审核通过" {
		return errors.New("无效参数 remark")
	}
	subscription := KcSubscription{Id: id}
	o := orm.NewOrm()
	o.Begin()
	if err := o.Read(&subscription); err != nil {
		o.Rollback()
		return err
	}
	if err := o.Read(&subscription, "order"); err != nil {
		o.Rollback()
		return err
	}
	if subscription.Status != 1 {
		o.Rollback()
		return errors.New("状态已改变")
	}
	subscription.AuthTime = time.Now()
	subscription.Status = status
	subscription.Remark = remark
	if _, err := o.Update(&subscription, "status", "remark", "auth_time"); err != nil {
		o.Rollback()
		return err
	}
	submissionStatus := 2
	if status == 0 {
		submissionStatus = 1
	}
	if _, err := o.QueryTable("kc_subscription_submission").Filter("order", subscription.Order).Filter("status", 0).Update(orm.Params{
		"AuthTime": time.Now(),
		"Status": submissionStatus,
		"Remark": remark,
	}); err != nil {
		o.Rollback()
		return err
	}
	if status == 2 {
		if err := AddAmount(o, subscription.Uid, subscription.CurAmount, subscription.Currency, "subscription"); err != nil {
			o.Rollback()
			beego.Error(err)
			return result.ErrCode(100102)
		}
		u := GetUserById(subscription.Uid)
		//认购直接锁仓
		switch utils.Period {
		case utils.Footstone:
			expire := utils.FootstoneExpireDate
			if err := handleLocked(o, u, subscription.Currency, subscription.CurAmount, 0, expire); err != nil {
				o.Rollback()
				return err
			}
		case utils.Angel:
			expire := time.Now().AddDate(0, 0, 60)
			if err := handleLocked(o, u, subscription.Currency, subscription.CurAmount, 1, expire); err != nil {
				o.Rollback()
				return err
			}
		}
	}
	o.Commit()
	if status == 2 {
		message := fmt.Sprintf("获得 %f %s", subscription.CurAmount, subscription.Currency)
		messageEn := fmt.Sprintf("Get %f %s", subscription.CurAmount, subscription.Currency)
		messageKo := fmt.Sprintf("얻다 %f %s", subscription.CurAmount, subscription.Currency)
		messageJp := fmt.Sprintf("獲得 %f %s", subscription.CurAmount, subscription.Currency)
		SetMessage(subscription.Uid, message, messageEn, messageKo, messageJp, "")
	}
	return nil
}
