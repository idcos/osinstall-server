package util

import (
	//"strings"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

func CreateNewMacAddress() string {
	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%s", time.Now().UnixNano())))
	cipherStr := h.Sum(nil)
	md5 := fmt.Sprintf("%s", hex.EncodeToString(cipherStr))
	var mac string
	mac = "52:54:00"
	mac += ":" + md5[0:2]
	mac += ":" + md5[2:4]
	mac += ":" + md5[4:6]
	return mac
}
