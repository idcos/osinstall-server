package main

import (
	"flag"
	"fmt"
	"server/osinstallserver/util"
)

func main() {
	var str string
	var publicKey string
	flag.StringVar(&str, "s", "", "需要加密的字符串")
	flag.StringVar(&publicKey, "p", "", "公钥地址")
	flag.Parse()
	if str == "" {
		fmt.Println("请设置需要加密的字符串参数，格式如: ./osinstall-encrypt-generator -p=\"./rsa/public.pem\" -s=\"user:password@tcp(localhost:3306)/idcos-osinstall?charset=utf8&parseTime=True&loc=Local\"")
		return
	}

	if publicKey == "" {
		fmt.Println("请设置公钥参数，格式如: ./osinstall-encrypt-generator -p=\"./rsa/public.pem\" -s=\"user:password@tcp(localhost:3306)/idcos-osinstall?charset=utf8&parseTime=True&loc=Local\"")
		return
	}

	str, err := util.RSAEncrypt(publicKey, str)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("生成成功，请将以下配置信息复制到/etc/osinstall-server/osinstall-server.conf里的\"repo\": {}节点里，并去除原有的\"connection\"子节点:")
	fmt.Println("\"connectionIsCrypted\":\"" + str + "\"")
}
