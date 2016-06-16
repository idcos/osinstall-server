package util

import (
	"io/ioutil"
	"os"
	"strings"
)

func CreatePxeFile(dir string, file string, content string) error {
	if !FileExist(dir) {
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			return err
		}
	}

	//文件已存在,先删除
	if FileExist(dir + "/" + file) {
		err := os.Remove(dir + "/" + file)
		if err != nil {
			return err
		}
	}

	bytes := []byte(content)
	err := ioutil.WriteFile(dir+"/"+file, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func FileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func ReadFile(path string) (string, error) {
	fi, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer fi.Close()
	bytes, err := ioutil.ReadAll(fi)
	return string(bytes), err
}

func GetPxeFileNameByMac(mac string) string {
	mac = strings.Replace(mac, ":", "-", -1)
	return "01-" + mac
}
