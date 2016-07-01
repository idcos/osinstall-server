package util

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"io/ioutil"
	"middleware"
	"os"
	"os/exec"
)

// execScript 执行脚本
func ExecScript(script string) ([]byte, error) {

	// 生成临时文件
	file, err := ioutil.TempFile("", "tmp-script")
	if err != nil {
		return nil, err
	}
	defer os.Remove(file.Name())
	defer file.Close()

	if _, err = file.WriteString(script); err != nil {
		return nil, err
	}
	file.Close()

	var cmd = exec.Command("/bin/bash", file.Name())
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	err = cmd.Wait()
	return output.Bytes(), err
}

func IconvFile(filename string, fromCode string, toCode string, ctx context.Context) error {
	if filename == "" || fromCode == "" || toCode == "" {
		return errors.New("参数错误")
	}

	logger, ok := middleware.LoggerFromContext(ctx)
	if !ok {
		return errors.New("内部服务器错误")
	}

	if !FileExist(filename) {
		return errors.New(filename + "文件不存在")
	}
	newFilename := filename + "_" + toCode
	cmd := fmt.Sprintf("iconv -t %s -f %s -c  %s > %s", toCode, fromCode, filename, newFilename)
	logger.Debugf("run script: %s", cmd)
	bytes, err := ExecScript(cmd)
	logger.Debugf("run result: %s", string(bytes))
	if err != nil {
		logger.Errorf("run error:%s", err.Error())
		return err
	}

	//delete old file
	errRemove := os.Remove(filename)
	if errRemove != nil {
		return errRemove
	}

	//rename
	errRename := os.Rename(newFilename, filename)
	if errRename != nil {
		return errRename
	}
	return nil
}
