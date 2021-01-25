package cryptutils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"strings"
)

// HMacSHa256 hmac sha256
func HMacSHa256(key, data string) (string, error) {
	h := hmac.New(sha256.New, []byte(key))
	_, err := io.WriteString(h, data)
	if err != nil {
		return "", err
	}
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil))), nil
}
