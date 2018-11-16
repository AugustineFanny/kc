package job

// import (
// 	"github.com/astaxie/beego/httplib"
// 	"fmt"
// 	"time"
// 	"errors"
// 	"github.com/spf13/cast"
// 	"github.com/satori/go.uuid"
// 	"github.com/gin-gonic/gin/json"
// 	"strings"
// 	"crypto/sha256"
// 	"encoding/hex"
// 	"crypto/hmac"
// 	"sync"
// 	"math/rand"
// 	"github.com/astaxie/beego"
// )

// const (
// 	API_HOST = "api.btcdo.com"
// )

// var (
// 	API_KEY = ""
// 	API_SECRET = ""

// 	Symbol = "ETH_USDT"
// 	Amount = 0.1
// 	Floatation = 0.5
// 	SellMinAmount = 0.0
// 	BuyMinAmount = 0.0
// 	Rest = 10000
// 	StartFlag = false
// 	TradeLogs = tradeLogs{}
// 	mu sync.Mutex
// )

// type tradeLogs struct {
// 	logs     []*tradeLogStruct
// }

// func (t *tradeLogs) Push(log string, value ...interface{}) {
// 	mu.Lock()
// 	defer mu.Unlock()
// 	if len(t.logs) >= 50 {
// 		t.logs = append(t.logs[1:], new(tradeLogStruct).Sprintf(log, value...))
// 	} else {
// 		t.logs = append(t.logs, new(tradeLogStruct).Sprintf(log, value...))
// 	}
// }

// func (t *tradeLogs) Pop() (res []*tradeLogStruct) {
// 	mu.Lock()
// 	defer mu.Unlock()
// 	res = t.logs
// 	t.logs = []*tradeLogStruct{}
// 	return
// }

// type tradeLogStruct struct {
// 	Time     string      `json:"time"`
// 	Log      string      `json:"log"`
// }

// func (t *tradeLogStruct) Sprintf(log string, value ...interface{}) *tradeLogStruct {
// 	t.Time = time.Now().Format("15:04:05.000")
// 	t.Log = fmt.Sprintf(log, value...)
// 	return t
// }

// type Depth struct {
// 	Price    float64
// 	Amount   float64
// }

// type Resp struct {
// 	Symbol    string
// 	BuyOrders []Depth
// 	SellOrders []Depth
// }

// func getDepth(symbol string) (*Depth, *Depth, error) {
// 	var resp Resp
// 	server := httplib.Get("https://api.btcdo.com/v1/market/depth/" + symbol).SetTimeout(10*time.Second, 10*time.Second)
// 	if err := server.ToJSON(&resp); err != nil {
// 		return nil, nil, err
// 	}
// 	if len(resp.BuyOrders) == 0 || len(resp.SellOrders) == 0 {
// 		return nil, nil, errors.New("err: 0000000000000000000000")
// 	}
// 	return &resp.SellOrders[0], &resp.BuyOrders[0], nil
// }

// func getOrders() {
// 	timestamp := cast.ToString(time.Now().Unix() * 1000)
// 	payloadList := []string{
// 		"GET",
// 		API_HOST,
// 		"/v1/trade/orders",
// 		"",
// 		"API-KEY: " + API_KEY,
// 		"API-SIGNATURE-METHOD: HmacSHA256",
// 		"API-SIGNATURE-VERSION: 1",
// 		"API-TIMESTAMP: " + timestamp,
// 		"",
// 	}
// 	payloadStr := strings.Join(payloadList, "\n")
// 	// signature
// 	mac := hmac.New(sha256.New, []byte(API_SECRET))
// 	mac.Write([]byte(payloadStr))
// 	md := mac.Sum(nil)
// 	mdStr := hex.EncodeToString(md)

// 	server := httplib.Get("https://api.btcdo.com/v1/trade/orders")
// 	server.Header("API-KEY", API_KEY)
// 	server.Header("API-SIGNATURE-METHOD", "HmacSHA256")
// 	server.Header("API-SIGNATURE-VERSION", "1")
// 	server.Header("API-TIMESTAMP", timestamp)
// 	server.Header("API-SIGNATURE", mdStr)
// 	server.Header("Content-Type", "application/json")
// 	return
// }

// type Order struct {
// 	Amount         float64    `json:"amount"`
// 	CustomFeatures int64      `json:"customFeatures"`
// 	OrderType      string     `json:"orderType"`
// 	Price          float64    `json:"price"`
// 	Symbol         string     `json:"symbol"`
// }

// func createOrder(symbol string, price, amount float64, orderType string) (string, error) {
// 	order := Order{Amount: amount, CustomFeatures: 65536, OrderType: orderType, Price: price, Symbol: symbol}
// 	str, err := json.Marshal(order)
// 	if err != nil {
// 		return "", err
// 	}

// 	timestamp := cast.ToString(time.Now().Unix() * 1000)
// 	u, err := uuid.NewV4()
// 	if err != nil {
// 		TradeLogs.Push("uuid error: %s", err)
// 		beego.Error(err)
// 	}
// 	uniqueId := hex.EncodeToString(u.Bytes())
// 	payloadList := []string{
// 		"POST",
// 		API_HOST,
// 		"/v1/trade/orders",
// 		"",
// 		"API-KEY: " + API_KEY,
// 		"API-SIGNATURE-METHOD: HmacSHA256",
// 		"API-SIGNATURE-VERSION: 1",
// 		"API-TIMESTAMP: " + timestamp,
// 		"API-UNIQUE-ID: " + uniqueId,
// 		string(str),
// 	}
// 	payloadStr := strings.Join(payloadList, "\n")

// 	// signature
// 	mac := hmac.New(sha256.New, []byte(API_SECRET))
// 	mac.Write([]byte(payloadStr))
// 	md := mac.Sum(nil)
// 	mdStr := hex.EncodeToString(md)

// 	server := httplib.Post("https://api.btcdo.com/v1/trade/orders")
// 	server.Header("API-KEY", API_KEY)
// 	server.Header("API-SIGNATURE-METHOD", "HmacSHA256")
// 	server.Header("API-SIGNATURE-VERSION", "1")
// 	server.Header("API-TIMESTAMP", timestamp)
// 	server.Header("API-UNIQUE-ID", uniqueId)
// 	server.Header("API-SIGNATURE", mdStr)
// 	server.Body(str)
// 	server.Header("Content-Type", "application/json")
// 	return server.String()
// }

// func sell(symbol string, price, amount float64, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	str, err := createOrder(symbol, price, amount, "SELL_LIMIT")
// 	if err != nil {
// 		TradeLogs.Push(err.Error())
// 		beego.Error(err)
// 	} else {
// 		TradeLogs.Push(str)
// 		beego.Error(str)
// 	}
// 	TradeLogs.Push("sell %s %f %f\n", Symbol, price, amount)
// }

// func buy(symbol string, price, amount float64, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	str, err := createOrder(symbol, price, amount, "BUY_LIMIT")
// 	if err != nil {
// 		TradeLogs.Push(err.Error())
// 		beego.Error(err)
// 	} else {
// 		TradeLogs.Push(str)
// 		beego.Error(str)
// 	}
// 	TradeLogs.Push("buy %s %f %f\n", Symbol, price, amount)
// }

// func handle() {
// 	symbol := Symbol
// 	amount := Amount
// 	floatation := Floatation
// 	//sellMinAmount := SellMinAmount
// 	//buyMinAmount := BuyMinAmount
// 	var wg sync.WaitGroup
// 	//sellOrder, buyOrder, err := getDepth(symbol)
// 	//if err != nil {
// 	//	TradeLogs.Push(err.Error())
// 	//	beego.Error(err)
// 	//	return
// 	//}
// 	//if sellOrder.Amount >= sellMinAmount && buyOrder.Amount >= buyMinAmount {
// 	//	destPrice := buyOrder.Price + (sellOrder.Price - buyOrder.Price) * floatation
// 	destPrice := floatation
// 	TradeLogs.Push("%f", destPrice)
// 	beego.Error(destPrice)
// 	if destPrice > 0 {
// 		wg.Add(2)
// 		go sell(symbol, destPrice, amount, &wg)
// 		go buy(symbol, destPrice, amount, &wg)
// 		wg.Wait()
// 	}
// 	//} else {
// 	//	TradeLogs.Push("empty")
// 	//	beego.Error("empty")
// 	//}
// }

// func BtcdoJob() {
// 	rand.Seed(time.Now().Unix())
// 	for {
// 		if StartFlag {
// 			handle()
// 		}
// 		time.Sleep(time.Duration(Rest) * time.Millisecond)
// 	}
// }
