package exec

import (
	"idcos.io/osinstall/logger"
	"idcos.io/osinstall/model"
	"idcos.io/osinstall/utils"
	"io/ioutil"
	"net/http"
	"strings"
)

func TaskExec(logger logger.Logger, reqParam *model.ConfJobIPExecParam, URL string) {
	client := &http.Client{}

	jsonParam := utils.ToJsonString(reqParam)

	logger.Debugf("[TaskExec]开始调用act2接口，url:%s, param:%s", URL, jsonParam)

	req, err := http.NewRequest("POST", URL, strings.NewReader(utils.ToJsonString(reqParam)))
	if err != nil {
		logger.Errorf("[TaskExec]调用act2接口异常， error: %s, url:%s, param:%s", err.Error(), URL, utils.ToJsonString(reqParam))
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	if resp == nil {
		logger.Errorf("[TaskExec]调用act2接口结束,返回值为nil")
		return
	}

	defer resp.Body.Close()

	logger.Debugf("do done")

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("解析response body error: %s, body: %s", err.Error(), body)
	}
	logger.Debugf("[TaskExec]调用act2接口结束,返回值为 %s", body)

}
