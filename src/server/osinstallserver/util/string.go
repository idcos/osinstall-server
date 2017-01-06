package util

import (
	//"strings"
	"errors"
)

func SubString(str string, begin int, length int) string {
	// 将字符串的转换成[]rune
	rs := []rune(str)
	lth := len(rs)

	// 简单的越界判断
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}

	end := begin + length
	if end > lth {
		end = lth
	}

	// 返回子串
	return string(rs[begin:end])
}

func CutArray(arr []string, length int) ([]string, []string, error) {
	var newArr []string
	var lostArr []string
	if len(arr) < length {
		return newArr, lostArr, errors.New("长度不够切分")
	}
	for i, value := range arr {
		if (i + 1) <= length {
			newArr = append(newArr, value)
		} else {
			lostArr = append(lostArr, value)
		}
	}
	return newArr, lostArr, nil
}

func IsInArray(str string, arr []string) bool {
	for _, value := range arr {
		if str == value {
			return true
		}
	}
	return false
}
