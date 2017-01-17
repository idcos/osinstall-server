package utils

import "github.com/axgle/mahonia"

// GBK2UTF8 将GBK字符串转换成UTF8字符串
func GBK2UTF8(strGBK string) (strUTF8 string) {
	return mahonia.NewDecoder("gbk").ConvertString(strGBK)
}

// UTF82GBK 将UTF8字符串转换成GBK字符串
func UTF82GBK(strUTF8 string) (strGBK string) {
	return mahonia.NewEncoder("gbk").ConvertString(strUTF8)
}
