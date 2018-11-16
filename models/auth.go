package models

import (
	"kuangchi_backend/result"
	"kuangchi_backend/utils"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"time"
	"encoding/base64"
	"strings"
)

func IsUsableEmail(email string) bool {
	o := orm.NewOrm()
	exist := o.QueryTable("user").Filter("email", email).Exist()
	if exist {
		return false
	}
	return true
}

func IsUsableMobile(mobile string) bool {
	o := orm.NewOrm()
	exist := o.QueryTable("user").Filter("mobile", mobile).Exist()
	if exist {
		return false
	}
	return true
}

func IsUsableUsername(username string) bool {
	o := orm.NewOrm()
	exist := o.QueryTable("user").Filter("username", username).Exist()
	if exist {
		return false
	}
	return true
}

func validUsername(username string) (err error) {
	valid := validation.Validation{}
	if v := valid.AlphaDash(username, "username"); !v.Ok {
		return result.ErrCode(100314)
	}
	if v := valid.MinSize(username, 4, "username"); !v.Ok {
		return result.ErrCode(100314)
	}
	if v := valid.MaxSize(username, 30, "username"); !v.Ok {
		return result.ErrCode(100314)
	}
	v := username[0]
	if ('Z' < v || v < 'A') && ('z' < v || v < 'a') {
		return result.ErrCode(100314)
	}
	return nil
}

func addInvite(u *User, inviteCode string) {
	o := orm.NewOrm()
	inviter := GetUserByInviteCode(inviteCode)
	if inviter != nil {
		//UserRecharge(inviter.Id, "BIT", 20)
		u.InviterId = inviter.Id
		u.Parents = inviter.DoAsParents()
		if _, err := o.Update(u, "InviterId", "Parents"); err != nil {
			beego.Error(err)
			return
		}
		invite := Invite{Uid: inviter.Id, Invite: u.Id}
		if _, err := o.Insert(&invite); err != nil {
			beego.Error(err)
			return
		}
		_, err := o.Raw("INSERT INTO invite(uid, invite, grade, create_time) SELECT uid, ?, grade + 1, NOW() FROM invite WHERE invite = ?", u.Id, inviter.Id).Exec()
		if err != nil {
			beego.Error(err)
		}
	}
}

func createUser(u *User, inviteCode string) (err error) {
	u.Salt = utils.RandomString(10)
	u.EncodePasswd()
	u.Language = "en"
	o := orm.NewOrm()
	u.InviteCode = utils.RandomCode(8)
	if o.QueryTable("user").Filter("InviteCode", u.InviteCode).Exist() {
		u.InviteCode = utils.RandomCode(8)
		if o.QueryTable("user").Filter("InviteCode", u.InviteCode).Exist() {
			u.InviteCode = utils.RandomCode(8)
		}
	}
	uid, err := o.Insert(u)
	if err != nil {
		beego.Error(err)
		return result.ErrCode(100102)
	}
	profile := Profile{Uid:int(uid)}
	_, err = o.Insert(&profile)
	if err != nil {
		beego.Error(err)
		return result.ErrCode(100102)
	}

	//UserRecharge(int(uid), "BIT", 30)
	if inviteCode != "" {
		addInvite(u, inviteCode)
	}

	return nil
}

func CreateUserByEmail(u *User, inviteCode string) error {
	valid := validation.Validation{}
	if v := valid.Email(u.Email, "email"); !v.Ok {
		return result.ErrCode(100301)
	}
	if !IsUsableEmail(u.Email) {
		return result.ErrCode(100305)
	}

	if err := validUsername(u.Username); err != nil {
		return err
	}

	if !IsUsableUsername(u.Username) {
		return result.ErrCode(100316)
	}
	if v := valid.MinSize(u.Password, 6, "password"); !v.Ok {
		return result.ErrCode(100306)
	}

	u.Status = 1
	if err := createUser(u, inviteCode); err != nil {
		return err
	}

	return nil
}

func CreateUserByMobile(u *User, inviteCode string) (err error) {
	valid := validation.Validation{}
	if !utils.ValidateMobile(u.Mobile) {
		return result.ErrCode(100311)
	}
	if !IsUsableMobile(u.Mobile) {
		return result.ErrCode(100313)
	}
	if err := validUsername(u.Username); err != nil {
		return err
	}

	if !IsUsableUsername(u.Username) {
		return result.ErrCode(100316)
	}

	if v := valid.MinSize(u.Password, 6, "password"); !v.Ok {
		return result.ErrCode(100306)
	}

	u.Status = 1
	if err := createUser(u, inviteCode); err != nil {
		return err
	}
	return nil
}

func Activate(key string) (err error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return result.ErrMsg("invalid link")
	}
	s := strings.Split(string(decodeBytes), "#")
	if len(s) != 2 {
		return result.ErrMsg("invalid link")
	}
	if !utils.CheckActivate(s[0], s[1]) {
		return result.ErrMsg("invalid link")
	}
	o := orm.NewOrm()
	user := User{Email: s[0]}
	if err := o.Read(&user, "Email"); err != nil {
		beego.Error(err)
		return result.ErrCode(100102)
	}
	user.Status = 1
	if _, err := o.Update(&user, "Status"); err != nil {
		beego.Error(err)
		return result.ErrCode(100102)
	}
	return nil
}

func ResetPasswordByEmail(email, newPassword string) (err error) {
	valid := validation.Validation{}
	if v := valid.MinSize(newPassword, 6, "newPassword"); !v.Ok {
		return result.ErrCode(100306)
	}
	o := orm.NewOrm()
	user := User{Email: email}
	if err := o.Read(&user, "Email"); err != nil {
		beego.Error(err)
		return result.ErrCode(100102)
	}

	if err := SetNewPassword(&user, newPassword); err != nil {
		return err
	}
	return nil
}

func ResetPasswordByMobile(mobile, newPassword string) (err error) {
	valid := validation.Validation{}
	if v := valid.MinSize(newPassword, 6, "newPassword"); !v.Ok {
		return result.ErrCode(100306)
	}
	o := orm.NewOrm()
	user := User{Mobile: mobile}
	if err := o.Read(&user, "Mobile"); err != nil {
		beego.Error(err)
		return result.ErrCode(100102)
	}

	if err := SetNewPassword(&user, newPassword); err != nil {
		return err
	}
	return nil
}

func UserSignInByEmail(email, password string) (*User, error) {
	valid := validation.Validation{}
	if v := valid.Email(email, "email"); !v.Ok {
		return nil, result.ErrCode(100301)
	}

	o := orm.NewOrm()
	user := User{Email: email}
	if err := o.Read(&user, "Email"); err != nil {
		return nil, result.ErrCode(100302)
	}
	if !user.ValidatePassword(password) {
		return nil, result.ErrCode(100303)
	}
	if user.Status == 0 {
		return nil, result.ErrCode(100317)
	}
	user.LastTime = time.Now()
	if _, err := o.Update(&user, "LastTime"); err != nil {
		beego.Error(err)
	}
	SendLoginNotifination(&user)
	return &user, nil
}

func UserSignInByMobile(mobile, password string) (*User, error) {
	valid := validation.Validation{}
	if !utils.ValidateMobile(mobile) {
		return nil, result.ErrCode(100311)
	}
	if v := valid.MinSize(password, 6, "password"); !v.Ok {
		return nil, result.ErrCode(100306)
	}

	o := orm.NewOrm()
	user := User{Mobile: mobile}
	if err := o.Read(&user, "Mobile"); err != nil {
		return nil, result.ErrCode(100312)
	}
	if !user.ValidatePassword(password) {
		return nil, result.ErrCode(100303)
	}
	if user.Status == 0 {
		return nil, result.ErrCode(100317)
	}
	user.LastTime = time.Now()
	if _, err := o.Update(&user, "LastTime"); err != nil {
		beego.Error(err)
	}
	SendLoginNotifination(&user)
	return &user, nil
}

func UserSignInByUsername(username, password string) (*User, error) {
	o := orm.NewOrm()
	user := User{Username: username}
	if err := o.Read(&user, "Username"); err != nil {
		return nil, result.ErrCode(100315)
	}
	if !user.ValidatePassword(password) {
		return nil, result.ErrCode(100303)
	}
	if user.Status == 0 {
		return nil, result.ErrCode(100317)
	}
	user.LastTime = time.Now()
	if _, err := o.Update(&user, "LastTime"); err != nil {
		beego.Error(err)
	}
	SendLoginNotifination(&user)
	return &user, nil
}
