package util

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

func EncodePassword(password string) (string, error) {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(password))
	cipherStr := md5Ctx.Sum(nil)
	result := strings.ToLower(hex.EncodeToString(cipherStr))
	return result, nil
}
