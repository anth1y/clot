package totp

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"math"
	"strconv"
	"strings"
	"time"
)

var digitsPower []int = []int{
	1,         // 0
	10,        // 1
	100,       // 2
	1000,      // 3
	10000,     // 4
	100000,    // 5
	1000000,   // 6
	10000000,  // 7
	100000000, // 8
}

type Token struct {
	Id              int
	Label           string
	Secret          string `json:"-"`
	EncryptedSecret string
	Digest          string
	Digits          int
	TimeStep        int
}

func (t *Token) Now() string {
	return t.At(time.Now())
}

func (t *Token) At(ts time.Time) string {
	period := (ts.Unix() / int64(t.TimeStep))
	msg := IntToByteArray(period)

	digest := t.HmacSHA(msg)

	offset := int(digest[19] & 0xf)
	code := int32(digest[offset]&0x7f)<<24 |
		int32(digest[offset+1]&0xff)<<16 |
		int32(digest[offset+2]&0xff)<<8 |
		int32(digest[offset+3]&0xff)

	otp := int64(code) % int64(math.Pow10(t.Digits))
	out := PadToken(t.Digits, strconv.FormatInt(otp, 10))

	return out
}

func (t *Token) Key() []byte {
	key, err := base32.StdEncoding.DecodeString(strings.ToUpper(t.Secret))
	if err != nil {
		panic(err)
	}

	return key
}

func (t *Token) HmacSHA(msg []byte) []byte {
	mac := hmac.New(sha1.New, t.Key())
	mac.Write(msg)

	return mac.Sum(nil)
}

func IntToByteArray(in int64) []byte {
	res := make([]byte, 8)

	i := len(res) - 1
	for i >= 0 {
		res[i] = byte(in & 0xff)
		i--
		in = in >> 8
	}

	return res
}

func PadToken(size int, text string) string {
	padSize := size - len(text)

	for i := 0; i < padSize; i++ {
		text = "0" + text
	}

	return text
}
