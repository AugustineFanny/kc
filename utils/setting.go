package utils

import (
	"encoding/json"
	"fmt"
	"kuangchi_backend/pkg/gomail.v2"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "github.com/astaxie/beego/cache/redis"
	"math/rand"
	"os"
	"path"
	"time"
	"io/ioutil"
	"strings"
	"strconv"
)

type resource struct {
	Value       string
	Name        string
}

type currency struct {
	Fee        []float64
	TradeFee   []float64   `json:"trade_fee"`
	Exchanges  []string
}

type Periods int

const (
	UpComping Periods = iota
	Footstone	//基石时期
	Angel		//天使时期
	Three		//之后
)

var (
	URL            string
	SecretKey      string
	FromBlock      string
	SuperMobile    string

	EmailServer    *gomail.Dialer
	CaptchaTimeout int64
	SmtpFromAddr   string

	RedisClient    RedisCache

	MediaPath      string
	BannerPath     string
	AvatarImgPath  string
	CardImgPath    string
	KycImgPath     string
	QrcodePath     string
	AppealPath     string
	SubmissionPath string

	PageSize       int64

	DefaultOptions *Options

	EtherServer    *Ether

	FishingCode    *fishingCode

	PaymentMethods []resource
	Countries      []resource
	Units          []resource
	PaymentMap     map[string]string
	CountryMap     map[string]string
	UnitMap        map[string]string

	RateBases      []string

	StartDate      time.Time
	FootstoneExpireDate time.Time
	Period		   Periods
)

const (
	BTCDecimal = 100000000
	ETHdecimal = 1000000000000000000
)

func newEmailServer() {
	smtpHost := beego.AppConfig.String("smtphost")
	smtpPort, _ := beego.AppConfig.Int("smtpport")
	smtpUser := beego.AppConfig.String("smtpuser")
	smtpPass := beego.AppConfig.String("smtppass")
	SmtpFromAddr = beego.AppConfig.String("smtpfromaddr")

	EmailServer = gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	var err error
	if CaptchaTimeout, err = beego.AppConfig.Int64("captchatimeout"); err != nil {
		panic("captchatimeout配置有误")
	}
}

func newRedisServer() {
	cacheConfig := beego.AppConfig.String("cacheconfig")
	cache := strings.Split(cacheConfig, "/")
	dbNum := 0
	if len(cache) > 1 {
		i, err := strconv.ParseInt(cache[1], 10, 64)
		if err != nil {
			panic(err)
		}
		dbNum = int(i)
	}
	RedisClient = RedisCache{conninfo: cache[0], dbNum: dbNum}
	RedisClient.connectInit()
	if _, err := RedisClient.Do("PING"); err != nil {
		panic(err)
	}
}

func initLog() {
	logPath := beego.AppConfig.String("logpath")
	logConfig := map[string]interface{}{
		"filename": logPath,
		"separate": []string{"error"},
	}
	b, _ := json.Marshal(logConfig)
	if err := logs.SetLogger(logs.AdapterMultiFile, string(b)); err != nil {
		panic(err)
	}
}

func initMedia() {
	// 上传文件存储位置
	MediaPath = beego.AppConfig.String("mediapath")
	if _, err := os.Stat(MediaPath); os.IsNotExist(err) {
		panic(fmt.Sprintf("mediapath路径不存在: %s", MediaPath))
	}
	//横幅
	//BannerPath = path.Join(MediaPath, "banner")
	//_, err := os.Stat(AvatarImgPath)
	//if os.IsNotExist(err) {
	//	os.MkdirAll(BannerPath, os.ModePerm)
	//}
	//头像
	//AvatarImgPath = path.Join(MediaPath, "avatar")
	//_, err = os.Stat(AvatarImgPath)
	//if os.IsNotExist(err) {
	//	os.MkdirAll(AvatarImgPath, os.ModePerm)
	//}
	//实名认证文件
	//CardImgPath = path.Join(MediaPath, "name_auth")
	//_, err = os.Stat(CardImgPath)
	//if os.IsNotExist(err) {
	//	os.MkdirAll(CardImgPath, os.ModePerm)
	//}
	//KYC文件
	//KycImgPath = path.Join(MediaPath, "kyc")
	//_, err = os.Stat(KycImgPath)
	//if os.IsNotExist(err) {
	//	os.MkdirAll(KycImgPath, os.ModePerm)
	//}
	//支付二维码
	//QrcodePath = path.Join(MediaPath, "qrcode")
	//_, err = os.Stat(QrcodePath)
	//if os.IsNotExist(err) {
	//	os.MkdirAll(QrcodePath, os.ModePerm)
	//}
	//申诉文件
	//AppealPath = path.Join(MediaPath, "appeal")
	//_, err = os.Stat(AppealPath)
	//if os.IsNotExist(err) {
	//	os.MkdirAll(AppealPath, os.ModePerm)
	//}
	//打币截图文件
	SubmissionPath = path.Join(MediaPath, "submission")
	_, err := os.Stat(SubmissionPath)
	if os.IsNotExist(err) {
		os.MkdirAll(SubmissionPath, os.ModePerm)
	}
}

func initResources() {
	buf, err := ioutil.ReadFile("resources.json")
	if err != nil {
		panic(err)
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(buf, &raw); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(raw["rate_bases"], &RateBases); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(raw["payment_methods"], &PaymentMethods); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(raw["countries"], &Countries); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(raw["units"], &Units); err != nil {
		panic(err)
	}
	PaymentMap = make(map[string]string)
	CountryMap = make(map[string]string)
	UnitMap = make(map[string]string)
	for _, method := range PaymentMethods {
		PaymentMap[method.Value] = method.Name
	}
	for _, country := range Countries {
		CountryMap[country.Value] = country.Name
	}
	for _, unit := range Units {
		UnitMap[unit.Value] = unit.Name
	}
}

func otherInit() {
	// 初始化随机数种子
	rand.Seed(time.Now().UnixNano())
	// 分页大小
	pageSize, err := beego.AppConfig.Int64("pagesize")
	if err != nil {
		panic("pagesize配置有误")
	}
	PageSize = pageSize

	DefaultOptions = NewOptions()

	FishingCode, err = NewFishingCode()
	if err != nil {
		panic(err)
	}

	URL = beego.AppConfig.String("url")
	SecretKey = beego.AppConfig.String("secretkey")
	FromBlock = beego.AppConfig.String("fromblock")
	SuperMobile = beego.AppConfig.String("supermobile")
	StartDate = time.Date(2018, 3, 9, 0, 0, 0, 0, time.UTC)
	FootstoneExpireDate = StartDate.AddDate(0, 0, 60 + 60)
	if time.Now().Sub(StartDate).Hours() < 0 {
		Period = UpComping
	} else if time.Now().Sub(StartDate).Hours() < 60 * 24 {
		Period = Footstone
	} else if time.Now().Sub(StartDate).Hours() < 180 * 24 {
		Period = Angel
	} else {
		Period = Three
	}
}

func init() {
	newEmailServer()
	beego.Info("newEmailServer over")
	newRedisServer()
	beego.Info("newRedisServer over")
	initLog()
	beego.Info("initLog over")
	initMedia()
	beego.Info("initMedia over")
	initResources()
	beego.Info("initResources over")
	otherInit()
	beego.Info("otherInit over")
}
