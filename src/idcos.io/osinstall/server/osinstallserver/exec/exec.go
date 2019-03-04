package exec

import (
	"encoding/json"
	"fmt"
	"idcos.io/osinstall/logger"
	"idcos.io/osinstall/model"
	"idcos.io/osinstall/utils"
	"io/ioutil"
	"net/http"
	"strings"
)

type Act2Resp struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Content string `json:"content"`
}

func TaskExec(logger logger.Logger, reqParam *model.ConfJobIPExecParam, URL string) Act2Resp {
	client := &http.Client{}

	jsonParam := utils.ToJsonString(reqParam)

	logger.Debugf("[TaskExec]开始调用act2接口，url:%s, param:%s", URL, jsonParam)

	req, err := http.NewRequest("POST", URL, strings.NewReader(utils.ToJsonString(reqParam)))
	if err != nil {
		str := fmt.Sprintf("[TaskExec]调用act2接口异常， error: %s, url:%s, param:%s", err.Error(), URL, utils.ToJsonString(reqParam))
		logger.Errorf(str)
		return Act2Resp{
			Status:  "failure",
			Message: str,
			Content: "",
		}
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	if resp == nil {
		str := "[TaskExec]调用act2接口结束,返回值为nil"
		logger.Errorf(str)
		logger.Errorf(str)
		return Act2Resp{
			Status:  "failure",
			Message: str,
			Content: "",
		}
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	logger.Debugf("[TaskExec]调用act2接口结束,返回值为 %s", body)

	if err != nil {
		str := fmt.Sprintf("解析response body error: %s, body: %s", err.Error(), body)
		logger.Errorf(str)
		return Act2Resp{
			Status:  "failure",
			Message: str,
			Content: "",
		}
	}

	var act2Resp Act2Resp
	if err := json.Unmarshal(body, &act2Resp); err != nil {
		str := fmt.Sprintf("解析response body error: %s, body: %s", err.Error(), body)
		logger.Errorf(str)
		return Act2Resp{
			Status:  "failure",
			Message: str,
			Content: "",
		}
	}
	return act2Resp
}
