package utils

//copy from https://github.com/balasanjay/totp
import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"hash"
	"net/url"
	"strconv"
	"time"
	"strings"

)

var digit_power = []int64{
	1,          // 0
	10,         // 1
	100,        // 2
	1000,       // 3
	10000,      // 4
	100000,     // 5
	1000000,    // 6
	10000000,   // 7
	100000000,  // 8
	1000000000, // 9
}

type Options struct {
	Time     func() time.Time
	Tries    []int64
	TimeStep time.Duration
	Digits   uint8
	Hash     func() hash.Hash
}

func (opt *Options) Url(label, secretkey string) (string, string) {
	u := &url.URL{
		Scheme: "otpauth",
		Host:   "totp",
		Path:   fmt.Sprintf("/%s", label),
	}

	secret := strings.TrimRight(base32.StdEncoding.EncodeToString([]byte(secretkey)), "=")

	params := url.Values{
		"secret": {secret},
	}

	u.RawQuery = params.Encode()
	return secret, u.String()
}

func (opt *Options) Authenticate(secretKey, userCode string) bool {

	if int(opt.Digits) != len(userCode) {
		return false
	}

	uc, err := strconv.ParseInt(userCode, 10, 64)
	if err != nil {
		return false
	}

	t := opt.Time().Unix() / int64(opt.TimeStep/time.Second)
	var tbuf [8]byte

	hm := hmac.New(opt.Hash, []byte(secretKey))
	var hashbuf []byte

	for i := 0; i < len(opt.Tries); i++ {
		b := t + opt.Tries[i]

		tbuf[0] = byte(b >> 56)
		tbuf[1] = byte(b >> 48)
		tbuf[2] = byte(b >> 40)
		tbuf[3] = byte(b >> 32)
		tbuf[4] = byte(b >> 24)
		tbuf[5] = byte(b >> 16)
		tbuf[6] = byte(b >> 8)
		tbuf[7] = byte(b)

		hm.Reset()
		hm.Write(tbuf[:])
		hashbuf = hm.Sum(hashbuf[:0])

		offset := hashbuf[len(hashbuf)-1] & 0xf
		truncatedHash := hashbuf[offset:]

		code := int64(truncatedHash[0])<<24 |
			int64(truncatedHash[1])<<16 |
			int64(truncatedHash[2])<<8 |
			int64(truncatedHash[3])

		code &= 0x7FFFFFFF
		code %= digit_power[len(userCode)]
		if code == uc {
			return true
		}
	}

	return false
}

func NewOptions() *Options {
	return &Options{
		Time:     time.Now,
		Tries:    []int64{0, -1},
		TimeStep: 30 * time.Second,
		Digits:   6,
		Hash:     sha1.New,
	}
}
