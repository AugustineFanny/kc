package utils

import (
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
)

type smsResp struct {
	RespCode string `json:"respCode"`
	RespMsg  string `json:"respMsg"`
}

/**
content :=  "注册验证码为【Eighteen】"
mobile :="13521290xxx,xxxxx"
*/
//func PushSms(content string, destMobiles string) error {
//
//	//# 短信提交地址，请联系管理员获取
//	url := "http://43.243.130.33:8860/sendSms"
//	//# 用户账号，必填
//	cust_code := "300161"
//	//# 用户密码，必填
//	cust_pwd := "FGHPF3JBGA"
//	//# 短信内容，必填
//	//# content = "{}于{} 收到扫码支付 {} 元【Eighteen】"
//	//# 接收号码，必填，同时发送给多个号码时,号码之间用英文半角逗号分隔
//	//# destMobiles = "{}".format(mobile)
//	//# 业务标识，选填，由客户自行填写不超过20位的数字
//	//uid := ""
//	//# 长号码，选填
//	//sp_code := ""
//	//# 是否需要状态报告
//	//need_report := "yes"
//	//# 数字签名，签名内容根据 “短信内容+客户密码”进行MD5编码后获得
//	sign := content + cust_pwd
//	//sign = sign.encode('utf-8')
//	//m = hashlib.md5()
//	//m.update(sign)
//
//	md5Ctx := md5.New()
//	md5Ctx.Write([]byte(sign))
//	cipherStr := md5Ctx.Sum(nil)
//
//	//sign = m.hexdigest()
//	sign = hex.EncodeToString(cipherStr)
//
//	sendst := `{"cust_code":"` + cust_code + `","sp_code":"","content":"` + content + `","destMobiles":"` + destMobiles + `","uid":"","need_report":"yes","sign":"` + sign + `"}`
//	beego.Debug(sendst)
//	//res=
//	//{"uid":"","status":"success","respCode":"0","respMsg":"提交成功！","totalChargeNum":1,"result":[{"msgid":"59106108291904934226","mobile":"13521290790","code":"0","msg":"提交成功.","chargeNum":1}]}
//	resp, err := http.Post(url, "application/json", strings.NewReader(sendst))
//	if err != nil {
//		beego.Error(err)
//		return err
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		beego.Error(err)
//		return err
//	}
//	var sms_resp smsResp
//	json.Unmarshal(body, &sms_resp)
//	if sms_resp.RespCode != "0" {
//		beego.Error(sms_resp.RespMsg)
//		return errors.New(sms_resp.RespMsg)
//	}
//	beego.Debug(sms_resp)
//	return nil
//}
func PushSms(content string, destMobiles string) error {
	url := "https://www.wondermary.com/api/sendsms.php?email=shaoyongvkbcel@gmail.com&key=631f850f40d413d08f4ef78251892818&recipient="
	url += destMobiles + "&message=" + content
	beego.Warn(url)
	resp, err := http.Get(url)
	if err != nil {
		beego.Error(err)
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		beego.Error(err)
		return err
	}
	//todo解析xml
	return nil
}
