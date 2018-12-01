package models

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/gob"
	"fmt"
	"time"
	"kuangchi_backend/pkg/x/crypto/pbkdf2"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strings"
)

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username" orm:"unique"` //用户名
	Avatar   string `json:"avatar"` //头像
	Email    string `json:"email" orm:"index"`  //邮箱
	CountryCode string `json:"country_code"`    //86:中国
	Mobile   string `json:"mobile" orm:"index"` //手机
	Password string `json:"-"`      //密码
	Salt     string `json:"-"`      //盐
	TfSecret string `json:"-"`      //两步验证密钥
	TfOpened bool	`json:"tf_opened"` //两步验证是否开启 0:关闭 1:开启
	TfLogin  bool   `json:"tf_login"` //登录是否需要两步验证 0:关闭 1:开启
	FundPassword string `json:"-"`  //资金密码
	FundSalt string `json:"-"`      //盐
	Fund     bool   `json:"fund"`   //资金密码是否已设置 0:未设置 1:已设置
	FishingCode string `json:"fishing_code"` //防钓鱼码
	Status   int    `json:"status"` //账户状态 0:未激活 1:已激活 2:冻结
	Role     int    `json:"role"`   //实名认证 0:未认证 1:认证中 2:认证不通过 3:认证通过
	Kyc      int    `json:"kyc"`    //kyc 0:未认证 1:认证中 2:认证不通过 3:认证通过
	Language string `json:"language"` //语言 zh_hans:简体 zh_hant:繁体 en:英文
	GroupId  int    `json:"group_id"`
	InviteCode string `json:"code" orm:"null;unique"`
	InviterId int   `json:"inviter" orm:"index"`
	Parents  string `json:"parents"`
	CreateTime time.Time `json:"create_time" orm:"auto_now_add;type(datetime)"`
	LastTime time.Time   `json:"last_time" orm:"auto_now_add;type(datetime)"`
}

func (u *User) EncodePasswd() {
	newPassword := pbkdf2.Key([]byte(u.Password), []byte(u.Salt), 10000, 50, sha256.New)
	u.Password = fmt.Sprintf("%x", newPassword)
}

func (u *User) ValidatePassword(passwd string) bool {
	newUser := User{Password: passwd, Salt: u.Salt}
	newUser.EncodePasswd()
	return subtle.ConstantTimeCompare([]byte(u.Password), []byte(newUser.Password)) == 1
}

func (u *User) EncodeFundPasswd() {
	newFundPassword := pbkdf2.Key([]byte(u.FundPassword), []byte(u.FundSalt), 10000, 50, sha256.New)
	u.FundPassword = fmt.Sprintf("%x", newFundPassword)
}

func (u *User) ValidateFundPassword(fundPasswd string) bool {
	newUser := User{FundPassword: fundPasswd, FundSalt: u.FundSalt}
	newUser.EncodeFundPasswd()
	return subtle.ConstantTimeCompare([]byte(u.FundPassword), []byte(newUser.FundPassword)) == 1
}

func (u *User) Verified(t string) bool {
	switch t {
	case "email" :
		return u.Email != ""
	case "mobile" :
		return u.Mobile != ""
	case "role" :
		return u.Role == 3
	case "kyc" :
		return u.Kyc == 3
	default:
		return false
	}
}

func (u *User) CompleteMobile() string {
	return u.CountryCode + u.Mobile
}

func (u *User) DoAsParents() string {
	if u.Parents == "" {
		return fmt.Sprintf("%d", u.Id)
	} else {
		return fmt.Sprintf("%s,%d", u.Parents, u.Id)
	}
}

type Invite struct {
	Id       int    `json:"id"`
	Uid      int    `json:"uid" orm:"index"`
	Invite   int    `json:"invite" orm:"index"`
	Grade    int    `json:"grade"`
	CreateTime time.Time `json:"create_time" orm:"auto_now_add;type(datetime)"`
}

type Admin struct {
	Id       int    `json:"id"`
	Username string `json:"username" orm:"unique"`
	Mobile   string `json:"mobile"`
	Password string `json:"-"`      //密码
	Salt     string `json:"-"`      //盐
	FundPassword string `json:"-"`  //资金密码
	FundSalt string `json:"-"`      //盐
	Super    int    `json:"super"`  //超级管理员
	Status   int    `json:"status"` //状态 0:启用 1:禁用
	CreateTime time.Time `json:"create_time" orm:"auto_now_add;type(datetime)"`
	LastTime time.Time   `json:"last_time" orm:"auto_now_add;type(datetime)"`
}

func (u *Admin) EncodePasswd() {
	newPassword := pbkdf2.Key([]byte(u.Password), []byte(u.Salt), 10000, 50, sha256.New)
	u.Password = fmt.Sprintf("%x", newPassword)
}

func (u *Admin) ValidatePassword(passwd string) bool {
	newAdmin := Admin{Password: passwd, Salt: u.Salt}
	newAdmin.EncodePasswd()
	return subtle.ConstantTimeCompare([]byte(u.Password), []byte(newAdmin.Password)) == 1
}

func (u *Admin) EncodeFundPasswd() {
	newFundPassword := pbkdf2.Key([]byte(u.FundPassword), []byte(u.FundSalt), 10000, 50, sha256.New)
	u.FundPassword = fmt.Sprintf("%x", newFundPassword)
}

func (u *Admin) ValidateFundPassword(fundPasswd string) bool {
	newAdmin := Admin{FundPassword: fundPasswd, FundSalt: u.FundSalt}
	newAdmin.EncodeFundPasswd()
	return subtle.ConstantTimeCompare([]byte(u.FundPassword), []byte(newAdmin.FundPassword)) == 1
}

type Profile struct {
	Id          int       `json:"-"`
	Uid         int       `json:"uid"         orm:"index"`
	FirstTrade  time.Time `json:"first_trade" orm:"null"` //首次交易时间
	TradeTimes  int       `json:"trade_times"`  //累计交易次数
	CompleteTimes int     `json:"complete_times"`  //累计完成次数
	TradeTotal  float64   `json:"trade_total" orm:"digits(20);decimals(10)"` //累计交易量，单位：btc
	AdTimes     int       `json:"ad_times"`  //作为广告方,累计交易次数
	AdCompleteTimes int   `json:"ad_complete_times"` //作为广告方,累计完成次数
	SellTimes   int       `json:"sell_times"`   //累计出售次数
	RefuseTimes int       `json:"refuse_times"` //累计拒绝次数
	AveragePass int       `json:"average_pass"` //平均放行时间，单位：分钟

	Trust       int       `json:"trust"`    //信任人数
	Trusted     int       `json:"trusted"`  //被信任次数
	Shield      int       `json:"shield"`   //屏蔽人数

	Wechat      bool      `json:"wechat"` //是否设置微信支付方式
	Alipay      bool      `json:"alipay"` //是否设置支付宝支付方式
	Bank        bool      `json:"bank"`   //是否设置银行转账支付方式
}

type Kyc struct {
	Id          int       `json:"id"`
	Uid         int       `json:"uid" orm:"index"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	Mobile      string    `json:"mobile"`
	Birthday    string    `json:"birthday"`
	Country     string    `json:"country"`
	Province    string    `json:"province"`
	City        string    `json:"city"`
	Street      string    `json:"street"`
	PostCode    string    `json:"post_code"`
	IdentityDoc string    `json:"identity_document"`
	FundsSource string    `json:"funds_source"`
	IdFront     string    `json:"photo_id_front"`
	IdBack      string    `json:"photo_id_back"`
	IdHold      string    `json:"photo_id_hold"`
	BankAccount string    `json:"bank_account_photo"`
	Passport    string    `json:"passport"`
	Utility     string    `json:"utility"`
	Driver      string    `json:"driver_license"`
	TaxBill     string    `json:"tax_bill"`
	BankStatements string `json:"bank_statements"`
	Other       string    `json:"other"`
	CreateTime time.Time `json:"create_time" orm:"auto_now_add;type(datetime)"`
	AuthTime   time.Time `json:"auth_time"   orm:"null;type(datetime)"`
	Status      int       `json:"status"      orm:"default(1)"` //KYC 1:认证中 2:认证不通过 3:认证通过
	Desc       string    `json:"desc"`
}

type KcMessage struct {
	Id          int       `json:"id"`
	Uid         int       `json:"uid"           orm:"index"`
	Message     string    `json:"message"`
	MessageEn   string    `json:"message_en"`
	MessageJp   string    `json:"message_jp"`
	MessageKo   string    `json:"message_ko"`
	MessageType string    `json:"message_type"`
	Extra       string    `json:"extra"`
	Readed      bool      `json:"readed"` //0：未读，1：已读
	CreateTime  time.Time `json:"create_time"   orm:"auto_now_add;type(datetime)"`
}

type RealName struct {
	Id         int       `json:"id"`
	Uid        int       `json:"uid"`
	Country    string    `json:"country"`
	CredentialType string `json:"credential_type"`
	Name       string    `json:"name"`
	StartDate  string    `json:"start_date"`
	EndDate    string    `json:"end_date"`
	Card       string    `json:"card"`
	CardFront  string    `json:"card_front"`
	CardBack   string    `json:"card_back"`
	CardHold   string    `json:"card_hold"`
	CreateTime time.Time `json:"create_time" orm:"auto_now_add;type(datetime)"`
	AuthTime   time.Time `json:"auth_time"   orm:"null;type(datetime)"`
	Status     int       `json:"status"      orm:"default(1)"` //实名认证 1:认证中 2:认证不通过 3:认证通过
	Desc       string    `json:"desc"`
}

type KcCurrency struct {
	Id          int		 `json:"id"`
	Currency    string   `json:"currency" orm:"unique"`
	Fee         float64  `json:"fee"`//提现费用
	FeeEth      int      `json:"fee_eth"` //0:不使用ETH 1:使用ETH作为手续费
	FeeTrade    float64  `json:"fee_trade"` //交易手续费
	MinLock     float64  `json:"min_lock"` //单次最小锁仓数量
	Token       int      `json:"token"`   //是否ERC20代币 0:不是 1:是
	Contract    string   `json:"contract_address"` //智能合约地址
	ConfirmNum  int      `json:"confirm_num"` //充值需确认数 只做显示用
	Recharge    int      `json:"recharge"`  //是否可以充值 0:启用 1:禁用
	Withdraw    int      `json:"withdraw"`  //是否可以提现 0:启用 1:禁用
	MinAmount   float64  `json:"min_amount"` //最小提现数量
	Exchanges   string   `json:"exchanges"` //支持查询价格的交易所，逗号分隔
	TransFlag   int      `json:"trans_flag"` //是否可以交易 0:可以 1:不可以
	Decimals    int      `json:"decimals"`  //ERC20 decimal
	BasePrice   float64  `json:"base_price"` //currency对rmb价格
	MInterestRate float64 `json:"mining_interest_rate"` //挖矿利率
	CInterestRate float64 `json:"competition_interest_rate"` //推广竞赛利率
	HodlLast    time.Time `json:"hodl_last" orm:"auto_now_add;type(datetime)"` //最近一次持有增值日期
	ShareLast   time.Time `json:"share_last" orm:"auto_now_add;type(datetime)"` //最近一次分享增值日期
}

func (cur KcCurrency) ExchangesList() []string {
	if cur.Exchanges == "" {
		return []string{}
	}
	return strings.Split(cur.Exchanges, ",")
}

type KcWallet struct {
	Id           int       `json:"-"`
	Uid          int       `json:"uid"           orm:"index"`
	Currency     string    `json:"currency"      orm:"index"`
	Amount       float64   `json:"amount"        orm:"digits(20);decimals(10)"`
	LockAmount   float64   `json:"lock_amount"   orm:"digits(20);decimals(10)"`
	MiningAmount float64   `json:"mining_amount" orm:"digits(20);decimals(10)"`
	UpdateTime   time.Time `json:"update_time"   orm:"auto_now;type(datetime)"`
	Share        int       `json:"share"` //0:未获得推广收益 1:已获
}

func (u *KcWallet) TableIndex() [][]string {
	return [][]string{
		[]string{"Uid", "Currency"},
	}
}

func (u *KcWallet) UsableAmount() float64 {
	if u.LockAmount < 0 {
		beego.Error(fmt.Sprintf("%d lock_amount < 0", u.Uid))
		return u.Amount
	}
	return u.Amount - u.LockAmount
}

type KcLocked struct {
	Id           int       `json:"id"`
	Uid          int       `json:"uid"          orm:"index"`
	Currency     string    `json:"currency"     orm:"index"`
	Amount       float64   `json:"amount"       orm:"digits(20);decimals(10)"`
	StartDate    time.Time `json:"start_date"   orm:"type(date)"`
	Status       int       `json:"status"` //0:锁仓 1:解锁
	Remark       string    `json:"remark"` //KcMining id
	Share        int       `json:"share"` //0:未获得推广收益 1:已获
	Class        int       `json:"class"` //0:认购 1:锁仓
	CreateTime   time.Time `json:"create_time" orm:"auto_now_add;type(datetime)"`
}

type KcMining struct {
	Id           int       `json:"id"`
	Uid          int       `json:"uid"        orm:"index"`
	Currency     string    `json:"currency"`
	Reward       float64   `json:"reward"` //奖励
	Reward1      float64   `json:"reward1"` // 个人推广奖励
	Reward2      float64   `json:"reward2"` // 竞赛推广奖励
	Mining       float64   `json:"mining"` //锁仓数量
	Rate         float64   `json:"rate"` //总算力
	Percentage   float64   `json:"percentage"` //比例
	InterestRate float64   `json:"interest_rate"` //利率
	Description  string    `json:"description"` //说明 mining:持有挖矿 share:分享挖矿
	Remark       string    `json:"remark"` //备注
	CreateTime   time.Time `json:"create_time" orm:"auto_now_add;type(datetime)"`
}

type KcUserAddress struct {
	Id         int		`json:"id"`
	Uid        int    	`json:"uid" orm:"index"`
	Currency   string 	`json:"currency" orm:"index"`
	Version    int      `json:"-"`
	Address    string 	`json:"address" orm:"unique"`
	AddressIndex string `json:"address_index"`
	AllAmount  float64 	`json:"all_amount" orm:"digits(20);decimals(10)"` //累计充值数量
	NowAmount  float64 	`json:"now_amount" orm:"digits(20);decimals(10)"` //地址当前余额
	UnConfirmedAmount float64 `json:"unconfirmed_amount" orm:"digits(20);decimals(10)"` //未确认金额
	UpdateTime time.Time `json:"update_time" orm:"null;type(datetime)"`
}

type KcUserAddressRecord struct {
	Id         int     `json:"id"`
	Uid        int     `json:"uid"`
	Currency   string  `json:"currency"`
	Hash       string  `json:"hash"`
	From       string  `json:"from"`
	To         string  `json:"to"`
	Amount     float64 `json:"amount" orm:"digits(20);decimals(10)"`
	FeeAmount  float64 `json:"fee" orm:"digits(20);decimals(10)"` //平台收取的手续费，非链上手续费
	FeeCurrency string `json:"fee_currency"`
	Direction  int     `json:"direction"` //0:充值 1:提现 2:站内转 3:锁仓转
	Status     int     `json:"status"` //0:确认中 1:异常 2:完成 3:待审核 4:审核失败
	Confirmations int `json:"confirmations"` //确认数
	Desc       string  `json:"desc"`
	Remark     string  `json:"remark"`
	CreateTime time.Time `json:"create_time" orm:"auto_now_add;type(datetime)"`
	CheckTime  time.Time `json:"check_time" orm:"type(datetime);null"`
}

func (u *KcUserAddressRecord) AllAmount() float64 {
	if u.FeeCurrency == u.Currency {
		return u.Amount + u.FeeAmount
	}
	return u.Amount
}

func (u *KcUserAddressRecord) TableIndex() [][]string {
	return [][]string{
		[]string{"Currency", "Hash"},
	}
}

type KcFundChange struct {
	Id         int       `json:"-"`
	Uid        int       `json:"uid"         orm:"index"`
	Currency   string    `json:"currency"    orm:"index"`
	Amount     float64   `json:"amount"      orm:"digits(20);decimals(10)"`
	Direction  int       `json:"direction"` //0:加 1:减
	Desc       string    `json:"desc"` //说明 deposit:充币 distribute:分发 instation:站内转 withdraw:提币 mining:持有挖矿 share:分享挖矿 locked:锁仓倍增 inlocked:锁仓转让
	Remark     string    `json:"remark" orm:"type(text)"`
	CreateTime time.Time `json:"create_time" orm:"auto_now_add;type(datetime)"`
}

type KcAddressPool struct {
	Id           int    `json:"-"`
	Address      string `json:"address" orm:"unique"`
	Currency     string `json:"currency" orm:"index"`
	AddressIndex string `json:"address_index"`
	Flag         int    `json:"flag" orm:"default(0)"`  //0:可用 1:已使用 2:保留
}

func (u *KcAddressPool) TableUnique() [][]string {
	return [][]string{
		[]string{"Currency", "AddressIndex"},
	}
}

type AdminOperationLog struct {
	Id         int        `json:"id"`
	Aid        int        `json:"aid"`
	UserName   string     `json:"username"`
	Api        string     `json:"api"`
	Args       string     `json:"args" orm:"type(text)"`
	Sheet      string     `json:"sheet"` //主要操作表
	RowId      int        `json:"column_id"` //操作行id
	Status     bool       `json:"status"` //0:失败 1:成功
	Remark     string     `json:"remark"`
	CreateTime time.Time  `json:"create_time" orm:"auto_now_add;type(datetime)"`
}

type Group struct {
	Id       int    `json:"id"`
	Name     string `json:"name" orm:"unique"`
}

type KcSubscription struct {
	Id          int        `json:"id"`
	Uid         int        `json:"uid"`
	Currency    string     `json:"currency"`
	Base        string     `json:"base"`
	BaseAmount  float64    `json:"base_amount"`
	CurAmount   float64    `json:"currency_amount"`
	Exchange    float64    `json:"exchange"` //兑换比例
	Status      int        `json:"status"` //0:未付款 1:已付款，待审核 2:完成 3:取消 4:超时取消
	Class       int        `json:"class"` //0:认购
	Remark      string     `json:"remark"`
	CreateTime  time.Time  `json:"create_time" orm:"auto_now_add;type(datetime)"`
	AuthTime    time.Time  `json:"auth_time"   orm:"null;type(datetime)"`
}

type KcSubscriptionSubmission struct {
	Id          int        `json:"id"`
	Order       string     `json:"order"`
	Txid        string     `json:"txid"`
	Screenshot  string     `json:"screenshot"`
	Status      int        `json:"status"`
	Remark      string     `json:"remark"`
	Warn1       string     `json:"warn1"`
	Warn2       string     `json:"warn2"`
	CreateTime  time.Time  `json:"create_time" orm:"auto_now_add;type(datetime)"`
	AuthTime    time.Time  `json:"auth_time"   orm:"null;type(datetime)"`
}

type KcSimulation struct {
	Id           int       `json:"id"`
	Currency     string    `json:"currency"`
	Username     string    `json:"username" orm:"unique"`
	MiningAmount float64   `json:"mining_amount" orm:"digits(20);decimals(10)"`
	Amount       float64   `json:"amount" orm:"digits(20);decimals(10)"`
	MPercentage  float64   `json:"m_percentage"`
	ShareAmount  float64   `json:"share_amount" orm:"digits(20);decimals(10)"`
	Reward1      float64   `json:"reward1"`
	Reward2      float64   `json:"reward2"`
	SPercentage  float64   `json:"s_percentage"`
	InviterId    string    `json:"inviter_id"`
	Parents      string    `json:"parents"`
	Remark       string    `json:"remark"`
}

type KcSUser struct {
	Id        int    `json:"id"`
	Username  string `json:"username" orm:"unique"` //用户名
	InviterId int   `json:"inviter_id" orm:"index"`
	Parents   string `json:"parents"`
	Inviter   string `json:"inviter"`
	Pnames    string `json:"parents_name"`
}

func (u *KcSUser) DoAsParents() string {
	if u.Parents == "" {
		return fmt.Sprintf("%d", u.Id)
	} else {
		return fmt.Sprintf("%s,%d", u.Parents, u.Id)
	}
}

func (u *KcSUser) DoAsParentsName() string {
	if u.Pnames == "" {
		return fmt.Sprintf("%s", u.Username)
	} else {
		return fmt.Sprintf("%s,%s", u.Pnames, u.Username)
	}
}

type KcSCurrency struct {
	Id          int		 `json:"id"`
	Currency    string   `json:"currency" orm:"unique"`
	MinLock     float64  `json:"min_lock"` //单次最小锁仓数量
	MInterestRate float64 `json:"mining_interest_rate"` //挖矿利率
	CInterestRate float64 `json:"competition_interest_rate"` //推广竞赛利率
}

type KcSActivity struct {
	Id           int       `json:"id"`
	Uid          int       `json:"uid"`
	Username     string    `json:"username"`
	InviterId    int       `json:"inviter" orm:"index"`
	Parents      string    `json:"parents"`
	Inviter      string    `json:"inviter"`
	Pnames       string    `json:"parents_name"`
	Subscription float64   `json:"subscription"`
	Lock         float64   `json:"lock"`
	In           float64   `json:"in"`
	InSource     string    `json:"in_source"`
	InSourceUid  int       `json:"in_source_id"`
	Out          float64   `json:"out"`
	OutDest      string    `json:"out_dest"`
	OutDestUid   int       `json:"out_dest_id"`
	Date         time.Time `json:"date"   orm:"type(date)"`
	ExpireDate   time.Time `json:"expire_date" orm:"type(date)"`
	ExpireDates  string    `json:"expire_dates" orm:"type(date)"`
}

func (u *KcSActivity) Mining() float64 {
	return u.Subscription + u.Lock
}

func (u *KcSActivity) handleExpireDates() {
	expireDatesList := []string{}
	for i := 0; i < 10; i++ {
		date := u.ExpireDate.AddDate(0, 0, i * 30)
		expireDatesList = append(expireDatesList, date.Format("2006-01-02"))
	}
	u.ExpireDates = strings.Join(expireDatesList, ",")
}

type KcSLocked struct {
	Id           int       `json:"id"`
	Uid          int       `json:"uid"          orm:"index"`
	Currency     string    `json:"currency"     orm:"index"`
	Amount       float64   `json:"amount"       orm:"digits(20);decimals(10)"`
	TotalAmount  float64   `json:"total_amount" orm:"digits(20);decimals(10)"`
	StartDate    time.Time `json:"start_date"   orm:"type(date)"`
	ExpireDate   time.Time `json:"expire_date"  orm:"type(date)"`
	CreateTime   time.Time `json:"create_time" orm:"auto_now_add;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(User), new(Invite), new(RealName), new(KcWallet), new(KcLocked),
		new(KcUserAddress), new(KcUserAddressRecord), new(KcAddressPool),
		new(AdminOperationLog), new(KcMessage), new(Profile), new(Kyc),
		new(Admin), new(KcCurrency), new(Group), new(KcFundChange),
		new(KcSubscription), new(KcSubscriptionSubmission), new(KcMining),
		new(KcSimulation), new(KcSUser), new(KcSCurrency), new(KcSActivity), new(KcSLocked),
	)

	user := beego.AppConfig.String("mysqluser")
	pass := beego.AppConfig.String("mysqlpass")
	urls := beego.AppConfig.String("mysqlurls")
	db := beego.AppConfig.String("mysqldb")

	mysqlDataSource := fmt.Sprintf(`%s:%s@(%s)/%s?charset=utf8&loc=Asia%%2FShanghai`, user, pass, urls, db)
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", mysqlDataSource)
	if beego.BConfig.RunMode == "dev" {
		orm.Debug = true
	}
	orm.RunSyncdb("default", false, true)
	/**
	因为 session 内部采用了 gob 来注册存储的对象，例如 struct，所以如果你采用了非 memory 的引擎，
	请自己在 main.go 的 init 里面注册需要保存的这些结构体，不然会引起应用重启之后出现无法解析的错误
	*/
	gob.Register(new(User))
	gob.Register(new(Admin))
}
