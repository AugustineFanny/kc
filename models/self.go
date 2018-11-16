package models

import (
	"kuangchi_backend/result"
	"kuangchi_backend/utils"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"strings"
	"strconv"
)

func GetUserByUsername(username string) *User {
	o := orm.NewOrm()
	user := User{Username: username}
	if err := o.Read(&user, "Username"); err != nil {
		return nil
	}
	return &user
}

func GetUserByEmail(email string) *User {
	o := orm.NewOrm()
	user := User{Email: email}
	if err := o.Read(&user, "Email"); err != nil {
		return nil
	}
	return &user
}

func GetUserByMobile(mobile string) *User {
	o := orm.NewOrm()
	user := User{Mobile: mobile}
	if err := o.Read(&user, "Mobile"); err != nil {
		return nil
	}
	return &user
}

func GetUserByIdent(ident string) *User {
	if ident == "" {
		return nil
	}
	if strings.Contains(ident, "@") {
		return GetUserByEmail(ident)
	} else if ident[0] >= '0' && ident[0] <= '9' {
		return GetUserByMobile(ident)
	} else {
		return GetUserByUsername(ident)
	}
	return nil
}

func GetUserById(id int) *User {
	o := orm.NewOrm()
	user := User{Id: id}
	if err := o.Read(&user, "Id"); err != nil {
		return nil
	}
	return &user
}

func GetUserByInviteCode(inviteCode string) *User {
	o := orm.NewOrm()
	user := User{InviteCode: inviteCode}
	if err := o.Read(&user, "InviteCode"); err != nil {
		return nil
	}
	return &user
}

func SetNewPassword(u *User, newPassword string) (err error) {
	u.Salt = utils.RandomString(10)
	u.Password = newPassword
	u.EncodePasswd()

	o := orm.NewOrm()
	_, err = o.Update(u, "Password", "Salt")
	if err != nil {
		beego.Error(err)
		return result.ErrCode(100102)
	}
	//冻结账户24小时
	utils.RedisClient.Set("blocked:" + u.Username, true, 24 * 60 * 60)
	return nil
}

func ChangeUserPassword(u *User, password string, newPassword string) (err error) {
	valid := validation.Validation{}
	if !u.ValidatePassword(password) {
		return result.ErrCode(100303)
	}
	if v := valid.MinSize(newPassword, 6, "newPassword"); !v.Ok {
		return result.ErrCode(100306)
	}

	if err := SetNewPassword(u, newPassword); err != nil {
		return err
	}
	return nil
}

func ExistRealName(uid int) bool {
	o := orm.NewOrm()
	exist := o.QueryTable("real_name").Filter("uid", uid).Filter("status__in", 1, 3).Exist()
	if exist {
		return true
	}
	return false
}

func GetRealName(u *User) *RealName {
	o := orm.NewOrm()
	realName := RealName{Uid:u.Id, Status: 3}
	if err := o.Read(&realName, "Uid", "Status"); err != nil {
		return nil
	}
	return &realName
}

func NameAuth(u *User, realName RealName) (err error) {
	realName.Uid = u.Id
	o := orm.NewOrm()
	o.Begin()
	user := User{Id: u.Id}
	if err := o.ReadForUpdate(&user); err != nil {
		o.Rollback()
		return result.ErrCode(100102)
	}
	if user.Role == 1 || user.Role == 3 {
		o.Rollback()
		return nil
	}
	_, err = o.Insert(&realName)
	if err != nil {
		o.Rollback()
		beego.Error(err)
		return result.ErrCode(100102)
	}
	u.Role = 1
	_, err = o.Update(u, "Role")
	if err != nil {
		o.Rollback()
		beego.Error(err)
		return result.ErrCode(100102)
	}
	o.Commit()
	return nil
}

func ExistKyc(uid int) bool {
	o := orm.NewOrm()
	exist := o.QueryTable("kyc").Filter("uid", uid).Filter("status__in", 1, 3).Exist()
	if exist {
		return true
	}
	return false
}

func GetKyc(u *User) *Kyc {
	o := orm.NewOrm()
	kyc := Kyc{Uid:u.Id, Status: 3}
	if err := o.Read(&kyc, "Uid", "Status"); err != nil {
		return nil
	}
	return &kyc
}

func SetTwoFactor(u *User, secretKey string, opened bool) bool {
	o := orm.NewOrm()
	u.TfSecret = secretKey
	u.TfOpened = opened
	_, err := o.Update(u, "TfSecret", "TfOpened")
	if err != nil {
		beego.Error(err)
		return false
	}
	return true
}

func SetFundPassword(u *User, password string) (err error) {
	block := false
	valid := validation.Validation{}
	if v := valid.MinSize(password, 6, "password"); !v.Ok {
		return result.ErrCode(100306)
	}

	u.FundSalt = utils.RandomString(10)
	u.FundPassword = password
	u.EncodeFundPasswd()
	if u.Fund == false {
		u.Fund = true
	} else {
		block = true
	}

	o := orm.NewOrm()
	_, err = o.Update(u, "FundPassword", "FundSalt", "Fund")
	if err != nil {
		beego.Error(err)
		return result.ErrCode(100102)
	}
	if block {
		//冻结账户24小时
		utils.RedisClient.Set("blocked:" + u.Username, true, 24 * 60 * 60)
	}
	return nil
}

func BindEmail(u *User, email string) (err error) {
	u.Email = email
	o := orm.NewOrm()
	_, err = o.Update(u)
	if err != nil {
		beego.Error(err)
		return result.ErrCode(100102)
	}
	return nil
}

func BindMobile(u *User, mobile, countryCode string) (err error) {
	u.Mobile = mobile
	u.CountryCode = countryCode
	o := orm.NewOrm()
	_, err = o.Update(u)
	if err != nil {
		beego.Error(err)
		return result.ErrCode(100102)
	}
	return nil
}

func GetMessages(u *User, pageNo int64) *utils.Page {
	o := orm.NewOrm()
	query := o.QueryTable("kc_message").Filter("uid", u.Id)
	cnt, _ := query.Count()
	page := utils.SetPage(pageNo, cnt, 10)
	var messages []*KcMessage
	num, err := query.Limit(page.PageSize, page.Offset).OrderBy("-id").All(&messages)
	if err != nil {
		beego.Error(err)
		return nil
	} else if num >= 0 {
		page.List = messages
		return page
	}
	return nil
}

func MessagesAllRead(u *User) {
	o := orm.NewOrm()
	o.QueryTable("kc_message").Filter("uid", u.Id).Update(orm.Params{
		"readed": 1,
	})
}

func GetMsgNum(u *User) int64 {
	o := orm.NewOrm()
	cnt, err := o.QueryTable("kc_message").Filter("uid", u.Id).Filter("readed", false).Count()
	if err != nil {
		return 0
	}
	return cnt
}

func MessageRead(u *User, messageId int) {
	o := orm.NewOrm()
	_, err := o.Raw("UPDATE message SET readed = 1 WHERE id = ? AND uid = ?", messageId, u.Id).Exec()
	if err != nil {
		beego.Error(err)
	}
}

func GetInfo(u *User) map[string]interface{} {
	o := orm.NewOrm()
	res := make(map[string]interface{})
	var profile Profile
	err := o.QueryTable("profile").Filter("uid", u.Id).One(&profile)
	if err != nil {
		beego.Error(err)
		return nil
	}
	res["role"] = u.Role
	res["avatar"] = u.Avatar
	res["username"] = u.Username
	res["trust"] = profile.Trust
	res["trusted"] = profile.Trusted
	res["create_time"] = u.CreateTime
	res["first_trade"] = profile.FirstTrade
	res["trade_times"] = profile.TradeTimes
	res["trade_total"] = profile.TradeTotal
	res["average_pass"] = profile.AveragePass
	res["wechat"] = profile.Wechat
	res["alipay"] = profile.Alipay
	res["bank"] = profile.Bank
	res["email"] = u.Email
	res["mobile"] = utils.ReplaceAsterisk(u.Mobile)
	res["language"] = u.Language
	res["name"] = ""
	if u.Role == 3 {
		var realName RealName
		err := o.QueryTable("real_name").Filter("uid", u.Id).Filter("status", 3).One(&realName)
		if err != nil {
			beego.Error(err)
		} else {
			res["name"] = realName.Name
		}
	}
	return res
}

func SetFishingCode(u *User, fishingCode string) error {
	valid := validation.Validation{}
	if v := valid.AlphaDash(fishingCode, "fishingCode"); !v.Ok {
		return result.ErrCode(100106)
	}
	o := orm.NewOrm()
	u.FishingCode = fishingCode
	_, err := o.Update(u, "FishingCode")
	if err != nil {
		beego.Error(err)
		return result.ErrCode(100102)
	}
	return nil
}

func SetKYC(u *User, kyc Kyc) (err error) {
	kyc.Uid = u.Id
	o := orm.NewOrm()
	o.Begin()
	_, err = o.Insert(&kyc)
	if err != nil {
		o.Rollback()
		beego.Error(err)
		return result.ErrCode(100102)
	}
	if u.Kyc == 0 || u.Kyc == 2 {
		u.Kyc = 1
		_, err = o.Update(u, "Kyc")
		if err != nil {
			o.Rollback()
			beego.Error(err)
			return result.ErrCode(100102)
		}
	}
	o.Commit()
	return nil
}

func SetAvatar(u *User) (err error) {
	o := orm.NewOrm()
	_, err = o.Update(u, "Avatar")
	if err != nil {
		beego.Error(err)
		return result.ErrCode(100102)
	}
	return nil
}

func GetInvitelink(u *User) map[string]interface{} {
	//o := orm.NewOrm()
	//num, _ := o.QueryTable("user").Filter("InviterId", u.Id).Count()
	return map[string]interface{}{
		"url": utils.URL + "/kuangfront/emailregister?i=" + u.InviteCode,
		"pc_url": utils.URL + "/kuangfront_pc/register?i=" + u.InviteCode,
		//"num": num,
	}
}

func GetInviteNum(u *User) map[string]int {
	var params []orm.Params
	res := map[string]int{"M1": 0, "M2": 0, "M3": 0, "M4": 0, "M5": 0, "M6": 0}
	sql := `SELECT grade, count(1) num FROM invite WHERE uid = ? GROUP BY grade;`
	o := orm.NewOrm()
	if _, err := o.Raw(sql, u.Id).Values(&params); err != nil {
		beego.Error(err)
		return res
	}
	for _, invite := range params {
		num, err := strconv.Atoi(invite["num"].(string))
		if err != nil {
			beego.Error(err)
		}
		switch invite["grade"].(string) {
		case "0": res["M1"] = num
		case "1": res["M2"] = num
		case "2": res["M3"] = num
		case "3": res["M4"] = num
		case "4": res["M5"] = num
		case "5": res["M6"] = num
		}
	}
	return res
}

type rcTree struct {
	Name     string    `json:"name"`
	IsLeaf    bool     `json:"isLeaf"`
}

func GetInvites(u *User, username string) []*rcTree {
	var list []orm.Params
	res := []*rcTree{}
	destUid := u.Id
	o := orm.NewOrm()
	if username != "" && username != u.Username {
		userDest := GetUserByUsername(username)
		if userDest == nil {
			return res
		}
		if o.QueryTable("invite").Filter("Uid", u.Id).Filter("Invite", userDest.Id).Exist() == false {
			return res
		}
		destUid = userDest.Id
	}
	sql := `SELECT ai.invite, count(bi.invite) count, username FROM invite ai
			LEFT JOIN invite bi ON ai.invite = bi.uid
			LEFT JOIN user ON ai.invite = user.id
			WHERE ai.uid = ? AND ai.grade = 0
			GROUP BY ai.invite`

	if _, err := o.Raw(sql, destUid).Values(&list); err != nil {
		beego.Error(err)
		return res
	}
	for _, data := range list {
		if data["username"] != nil {
			rc := &rcTree{data["username"].(string), true}
			if data["count"].(string) != "0" {
				rc.IsLeaf = false
			}
			res = append(res, rc)
		}
	}
	return res
}
