package util

import (
	//"strings"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math"
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

func FotmatNumberToMB(number float64, unit string) int {
	var result float64
	if unit == "KiB" {
		result = number / 1024
	} else if unit == "MiB" {
		result = number
	} else if unit == "GiB" {
		result = number * 1024
	} else if unit == "TiB" {
		result = number * 1024 * 1024
	}
	return int(math.Floor(result))
}

func FotmatNumberToGB(number float64, unit string) int {
	var result float64
	if unit == "KiB" {
		result = number / 1024 / 1024
	} else if unit == "MiB" {
		result = number / 1024
	} else if unit == "GiB" {
		result = number
	} else if unit == "TiB" {
		result = number * 1024
	}
	return int(math.Floor(result))
}
