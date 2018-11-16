package utils

import (
	"math/big"
	"fmt"
	"encoding/json"
	"math/rand"
	"github.com/astaxie/beego/httplib"
	"errors"
)

type Ether struct {
	Url        string
	Address    string
	PassPhrase string
}

func NewEther(url, address, passphrase string) *Ether {
	return &Ether{
		Url: url,
		Address: address,
		PassPhrase: passphrase,
	}
}

type Tx struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
}

type clientRequest struct {
	// A String containing the name of the method to be invoked.
	Method string `json:"method"`
	// Object to pass as request parameter to the method.
	Params []interface{} `json:"params"`
	// The request id. This can be of any type. It is used to match the
	// response with the request that it is replying to.
	Id uint64 `json:"id"`
}

// clientResponse represents a JSON-RPC response returned to a client.
type clientResponse struct {
	Result *json.RawMessage `json:"result"`
	Error  interface{}      `json:"error"`
	Id     uint64           `json:"id"`
}

func EncodeBig(bigint *big.Int) string {
	nbits := bigint.BitLen()
	if nbits == 0 {
		return "0x0"
	}
	return fmt.Sprintf("%#x", bigint)
}

func EncodeEther(value float64) string {
	x := big.NewInt(ETHdecimal/1000000)
	y := big.NewInt(int64(value * 1000000))
	z := new(big.Int).Mul(x, y)
	return EncodeBig(z)
}

// EncodeClientRequest encodes parameters for a JSON-RPC client request.
func EncodeClientRequest(method string, args []interface{}) ([]byte, error) {
	c := &clientRequest{
		Method: method,
		Params: args,
		Id:     uint64(rand.Int63()),
	}
	return json.Marshal(c)
}

func (e *Ether) SendEther(to string, value float64) (string, error) {
	fmt.Println(e.Address, e.PassPhrase)
	tx := Tx{e.Address, to, EncodeEther(value)}
	args := []interface{}{tx, e.PassPhrase}
	a, _ := EncodeClientRequest("personal_sendTransaction", args)
	req := httplib.Post(e.Url)
	req.Body(a)
	b, _ := req.Bytes()
	fmt.Println(string(b))
	var c clientResponse
	json.Unmarshal(b, &c)
	if c.Error != nil {
		return "", errors.New(fmt.Sprintf("%v", c.Error))
	}
	if c.Result == nil {
		return "", errors.New("连接eth节点异常")
	}
	hash := string(*c.Result)
	return hash[1:len(hash)-1], nil
}
