package route

import (
	//"encoding/base64"
	//"fmt"
	"github.com/AlexanderChen1989/go-json-rest/rest"
	"golang.org/x/net/context"
	"middleware"
	//"net/http"
	//"strings"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func GetDeviceInstallReport(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	where := "`status` = 'success'"
	count, err := repo.CountDeviceInstallReportByWhere(where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	where = "`status` = 'success'"
	productReport, err := repo.GetDeviceProductNameInstallReport(where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	where = "`status` = 'success'"
	companyReport, err := repo.GetDeviceCompanyNameInstallReport(where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	where = "`status` = 'success'"
	osReport, err := repo.GetDeviceOsNameInstallReport(where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	where = "`status` = 'success'"
	hardwareReport, err := repo.GetDeviceHardwareNameInstallReport(where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	result := make(map[string]interface{})
	result["Count"] = count
	result["ProductReport"] = productReport
	result["CompanyReport"] = companyReport
	result["HardwareReport"] = hardwareReport
	result["OsReport"] = osReport

	w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功", "Content": result})
}

func ReportDeviceInstallReport(ctx context.Context, w rest.ResponseWriter, r *rest.Request) {
	repo, ok := middleware.RepoFromContext(ctx)
	if !ok {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "内部服务器错误"})
		return
	}

	where := " id > 0 and `status` = 'success'"
	count, err := repo.CountDeviceInstallReportByWhere(where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	where = "`status` = 'success'"
	productReport, err := repo.GetDeviceProductNameInstallReport(where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	where = "`status` = 'success'"
	hardwareReport, err := repo.GetDeviceHardwareNameInstallReport(where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	where = "`status` = 'success'"
	osReport, err := repo.GetDeviceOsNameInstallReport(where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	where = "`status` = 'success'"
	systemReport, err := repo.GetDeviceSystemNameInstallReport(where)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	result := make(map[string]interface{})
	result["Count"] = count
	result["ProductReport"] = productReport
	result["HardwareReport"] = hardwareReport
	result["OsReport"] = osReport
	result["SystemReport"] = systemReport

	b, err := json.Marshal(result)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	body := bytes.NewBuffer([]byte(b))
	resp, err := http.Post("http://open.idcos.com/api/x86/report-install-info", "application/json;charset=utf-8", body)
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": "网络连接故障，操作失败！您可以通过线下方式分享给我们，谢谢！"})
		return
	}
	reportResult, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": err.Error()})
		return
	}

	type Response struct {
		Status  string `json:"Status"`
		Message string `json:"Message"`
	}
	var response Response
	errJson := json.Unmarshal(reportResult, &response)
	if errJson != nil {
		w.WriteJSON(map[string]interface{}{"Status": "error", "Message": errJson.Error()})
		return
	}
	if response.Status == "success" {
		w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作成功"})
	} else {
		w.WriteJSON(map[string]interface{}{"Status": "success", "Message": "操作失败！您可以通过线下方式分享给我们，谢谢！"})
	}
}
