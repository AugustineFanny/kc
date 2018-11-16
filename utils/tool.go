package utils

import (
	"math/rand"
	"strconv"

	"kuangchi_backend/pkg/gomail.v2"
	"github.com/astaxie/beego"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"math"
	"time"
	"image"
	_ "image/png"
	"image/draw"
	"os"
	"io/ioutil"
	"io"
	"image/png"
	"bytes"
	"fmt"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"github.com/dgrijalva/jwt-go"
	"regexp"
)

// copy from https://github.com/moby/moby/blob/master/pkg/stringutils/stringutils.go
func RandomString(n int) string {
	letters := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RandomCaptcha(n int) string {
	letters := []byte("0123456789")
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RandomCode(n int) string {
	letters := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func SendRawEmail(email, title, body string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "FADAX Notification<" + SmtpFromAddr + ">")
	m.SetHeader("To", email)
	m.SetHeader("Subject", title)
	m.SetBody("text/html", body)
	beego.Debug("发送邮件")
	if err := EmailServer.DialAndSend(m); err != nil {
		beego.Error(err)
	} else {
		beego.Debug("发送成功")
	}
}

const emailhtml = `
<table width="800" border="0" align="center" cellpadding="0" cellspacing="0" bgcolor="#ffffff" style="font-family:'Microsoft YaHei';"><tbody>
	<tr>
		<td>
			<table width="800" border="0" align="left" cellpadding="0" cellspacing="0" style=" border:1px solid #edecec; border-top:none; border-bottom:none; padding:0 20px;font-size:14px;color:#333333;">
				 <tbody>
				 	<tr> <td width="760" height="56" border="0" align="left" colspan="2" style=" font-size:16px;vertical-align:bottom;font-family:'Microsoft YaHei';">Dear %s<a></a>:</td> </tr>
				 	<tr> <td width="760" height="30" border="0" align="left" colspan="2">&nbsp;</td> </tr>
					%s
				 	<tr> <td width="720" height="32" colspan="2" style="padding-left:40px;">&nbsp;</td> </tr>
				 	<tr> <td width="720" height="32" colspan="2" style="padding-left:40px;">&nbsp;</td> </tr>
				<tr>
					<td width="720" height="14" colspan="2" style="padding-bottom:16px; border-bottom:1px dashed #e5e5e5;font-family:'Microsoft YaHei';">FADAX.io</td>
				</tr>
				<tr>
					<td width="720" height="14" colspan="2" style="padding:8px 0 28px;color:#999999; font-size:12px;font-family:'Microsoft YaHei';">System mail</td>
				</tr>
			</tbody></table>
		</td>
	</tr>
</tbody></table>
`

const contenthtml = `
	<tr> <td width="720" height="32" border="0" align="left" valign="middle" style=" width:720px; text-align:left;vertical-align:middle;line-height:32px;font-family:'Microsoft YaHei';">%s</td> </tr>
`

const linkhtml = `
	<tr> <td width="720" height="32" colspan="2" style="padding-left:40px;font-family:'Microsoft YaHei';"><br><a href="%s" target="_blank">%s </a></td> </tr>
`

func buildContent(username, body string, link ...string) string {
	content := fmt.Sprintf(contenthtml, body)
	if len(link) > 0 {
		content += fmt.Sprintf(linkhtml, link[0], link[0])
	}
	return fmt.Sprintf(emailhtml, username, content)
}

func SendContentEmail(email, title, body string) {
	SendRawEmail(email, title, buildContent(email, body))
}

func SendEmail(email, captcha string) {
	body := fmt.Sprintf("Verification Code：%s，%d minutes of validity.", captcha, CaptchaTimeout / 60)
	SendRawEmail(email, "FADAX", buildContent(email, body))
}

func CheckActivate(email, sum string) bool {
	hash := sha256.New()
	hash.Write([]byte(email + SecretKey))
	md := hash.Sum(nil)
	if hex.EncodeToString(md) == sum {
		return true
	}
	return false
}

func CheckResetPassword(email, expire, sum string) bool {
	now := time.Now().Unix()
	timestamp, err := strconv.ParseInt(expire, 10, 64)
	if err != nil {
		return false
	}
	if now > timestamp {
		return false
	}
	hash := sha256.New()
	hash.Write([]byte(email + expire + SecretKey))
	md := hash.Sum(nil)
	if hex.EncodeToString(md) != sum {
		return false
	}
	return true
}

func SendSMS(mobile, captcha string) bool {
	body := fmt.Sprintf("Verification Code：%s，%d minutes of validity.", captcha, CaptchaTimeout / 60)
	if err := PushSms(body, mobile); err != nil {
		return false
	}
	return true
}

func CaptchaSet(identity, captcha, captchaType string) bool {
	key := "captcha:" + captchaType + ":" + identity
	err := RedisClient.Set(key, captcha, CaptchaTimeout)
	if err != nil {
		beego.Error(err)
		return false
	}
	return true
}

func CaptchaGet(identity, captchaType string) string {
	key := "captcha:" + captchaType + ":" + identity
	captcha, err := RedisClient.GetString(key)
	if err != nil {
		return ""
	}
	return captcha
}

func CaptchaDel(identity, captchaType string) {
	key := "captcha:" + captchaType + ":" + identity
	RedisClient.Delete(key)
}

func CheckCard(card string) bool {
	if len(card) != 18 {
		return false
	}
	wi := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	a18 := [11]byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}
	res := 0
	arr := make([]int, 17)

	for index, value := range card[:17] {
		arr[index], _ = strconv.Atoi(string(value))
	}

	for i := 0; i < 17; i++ {
		res += arr[i] * wi[i]
	}

	res = res % 11
	if a18[res] == card[17] {
		return true
	}
	if a18[res] == 'X' && card[17] == 'x' {
		return true
	}
	return false
}

func TimestampString() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

type Page struct {
	PageNo     int64       `json:"page_no"`
	PageSize   int64       `json:"page_size"`
	TotalPage  int64       `json:"total_page"`
	TotalCount int64       `json:"total_count"`
	HasPre     bool        `json:"has_pre"`
	HasNext    bool        `json:"has_next"`
	PrePage    int64       `json:"pre_page"`
	NextPage   int64       `json:"next_page"`
	Offset     int64       `json:"-"`
	List       interface{} `json:"list"`
}

func SetPage(pageNo, totalCount int64, size ...int64) *Page {
	pageSize := PageSize
	if len(size) > 0 {
		pageSize = size[0]
	}
	totalPage := int64(math.Ceil(float64(totalCount) / float64(pageSize)))
	if totalPage == 0 {
		totalPage = 1
	}
	hasPre, hasNext, prePage, nextPage := true, true, int64(0), int64(0)
	if pageNo <= 1 {
		pageNo = 1
		hasPre = false
	}
	if pageNo >= totalPage {
		pageNo = totalPage
		hasNext = false
	}
	if hasPre {
		prePage = pageNo - 1
	}
	if hasNext {
		nextPage = pageNo + 1
	}
	offset := (pageNo - 1) * pageSize
	return &Page{pageNo, pageSize, totalPage, totalCount, hasPre, hasNext, prePage, nextPage, offset, nil}
}

type fishingCode struct {
	size       float64
	dpi        float64
	font       *truetype.Font
	background image.Image
}

func NewFishingCode() (*fishingCode, error) {
	fontBytes, err := ioutil.ReadFile("luxirr.ttf")
	if err != nil {
		return nil, err
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}
	bg, err := os.Open("background.png")
	if err != nil {
		return nil, err
	}
	origin, _, err := image.Decode(bg)
	if err != nil {
		return nil, err
	}
	return &fishingCode{8, 144, font, origin}, nil
}

func (this *fishingCode) WriteTo(code string, writer io.Writer) {
	fg := image.White
	rgba := image.NewNRGBA(this.background.Bounds())
	draw.Draw(rgba, rgba.Bounds(), this.background, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(this.dpi)
	c.SetFont(this.font)
	c.SetFontSize(this.size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)

	// Draw the text.
	pt := freetype.Pt(30, 15+int(c.PointToFixed(this.size)>>6))
	c.DrawString(code, pt)
	var buf bytes.Buffer
	if err := png.Encode(&buf, rgba); err != nil {
		beego.Error(err)
	}
	writer.Write(buf.Bytes())
}

func GetCountry(code string) string {
	return CountryMap[code]
}

func GetPaymentMethods(lists ...string) string {
	res := []string{}
	for _, method := range lists {
		if PaymentMap[method] != "" {
			res = append(res, method)
		}
	}
	return strings.Join(res, ",")
}

func GetUnit(code string) string {
	return UnitMap[code]
}

func Sub(x, y float64) float64 {
	tmp := (x * 1000000 - y * 1000000) / 1000000
	s := strconv.FormatFloat(tmp, 'f', 6, 64)
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func CreateToken(username string) string {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Unix() + 86400 * 7,
		Issuer:    username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		beego.Error(err)
	}
	return ss
}

func CheckToken(authorization string) string {
	if strings.HasPrefix(authorization, "Bearer ") {
		type MyCustomClaims struct {
			jwt.StandardClaims
		}
		// sample token is expired.  override time so it parses as valid

		token, err := jwt.ParseWithClaims(authorization[7:], &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})

		if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
			if time.Now().Unix() < claims.StandardClaims.ExpiresAt {
				return claims.StandardClaims.Issuer
			}
		} else {
			beego.Error(err)
		}
	}
	return ""
}

func ShowFloat(f float64, digit int) float64 {
	//小数点后2位或6位
	switch digit {
	case 2:
		return math.Floor(f*100) / 100
	case 6:
		return math.Floor(f * 1000000) / 1000000
	case 8:
		return math.Floor(f * 100000000) / 100000000
	}
	return math.Floor(f * 1000000) / 1000000
}

//四舍五入
func RoundFloat(f float64, digit int) float64 {
	switch digit {
	case 2:
		return math.Floor(f*100 + 0.5) / 100
	case 6:
		return math.Floor(f*1000000 + 0.5) / 1000000
	}
	return math.Floor(f + 0.5)
}

func RealTimePrice(exchange, currency, unit string, premium float64) float64 {
	exPrice, err := RedisClient.HgetFloat64(fmt.Sprintf("exchange:%s:%s", exchange, currency), unit)
	if err != nil {
		exPrice = 0
	}
	if premium == 0 {
		return ShowFloat(exPrice, 2)
	}
	premium = premium / 100
	return ShowFloat(exPrice * (1 + premium), 2)
}

const asterisks = "********************"

func ReplaceAsterisk(str string) string {
	length := len(str)
	if length <= 4 {
		return str
	}
	return str[:2] + asterisks[:length-4] + str[length-2:]

}

func PasswordStrength(password string) bool {
	r, _ := regexp.MatchString(`\d`, password)
	c, _ := regexp.MatchString(`[a-zA-Z]`, password)
	return r && c
}

func ValidateMobile(mobile string) bool {
	match, _ := regexp.MatchString(`^\d+$`, mobile)
	return match
}
