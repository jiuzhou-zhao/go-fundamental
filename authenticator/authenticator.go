package authenticator

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"net/url"
	"time"

	"github.com/golang/glog"
)

// CreateGoogleAuthQRCodeData 创建二维码
func CreateGoogleAuthQRCodeData(key string, account string, issuer string) string {
	qrCode := "otpauth://totp/%s:%s?secret=%s&issuer=%s"
	qr := fmt.Sprintf(qrCode,
		encode(issuer),
		encode(account),
		encode(key),
		encode(issuer))
	glog.Infof("[CreateGoogleAuthQRCodeData] qr: %s", qr)
	return qr
}

// MakeGoogleAuthenticator 获取 key 对应的验证码，时间戳需要指定
func MakeGoogleAuthenticator(key string, t int64) (string, error) {
	hs, e := hmacSha1(key, t/30)
	if e != nil {
		return "", e
	}
	snum := lastBit4byte(hs)
	d := snum % 1000000
	return fmt.Sprintf("%06d", d), nil
}

// MakeGoogleAuthenticatorForNow 获取 key 对应的验证码，时间戳是当前
func MakeGoogleAuthenticatorForNow(key string) (string, error) {
	return MakeGoogleAuthenticator(key, time.Now().Unix())
}

func encode(str string) string {
	t := url.URL{Path: str}
	return t.String()
}

func lastBit4byte(hmacSha1 []byte) int32 {
	if len(hmacSha1) != sha1.Size {
		return 0
	}
	offsetBits := int8(hmacSha1[len(hmacSha1)-1]) & 0x0f
	p := (int32(hmacSha1[offsetBits]) << 24) | (int32(hmacSha1[offsetBits+1]) << 16) | (int32(hmacSha1[offsetBits+2]) << 8) | (int32(hmacSha1[offsetBits+3]) << 0)
	return (p & 0x7fffffff)
}

func hmacSha1(key string, t int64) ([]byte, error) {
	encoding := base32.StdEncoding.WithPadding(base32.NoPadding)
	decodeKey, err := encoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	cData := make([]byte, 8)
	binary.BigEndian.PutUint64(cData, uint64(t))

	h1 := hmac.New(sha1.New, decodeKey)
	_, e := h1.Write(cData)
	if e != nil {
		return nil, e
	}
	return h1.Sum(nil), nil
}
