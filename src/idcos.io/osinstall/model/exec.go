package model

const (
	DefaultIDC      = ""
	DefaultEncoding = "utf-8"

	Win   = "windows"
	Linux = "linux"
	Aix   = "aix"

	Root = "/tmp/root"

	CallbackURL = "/api/osinstall/v1/task/callback"
)

// ConfJobIPExecParam 作业执行参数
type ConfJobIPExecParam struct {
	ExecHosts   []ExecHost `json:"execHosts" validate:"required"`
	ExecParam   ExecParam  `json:"execParam" validate:"required"`
	Provider    string     `json:"provider" validate:"required"` // provider可以为salt|puppet|openssh
	Callback    string     `json:"callback"`
	ExecuteID   string     `json:"executeId" validate:"required"`
	JobRecordID string     `json:"jobRecordId"`
}

//ExecHost
type ExecHost struct {
	HostIP   string `json:"hostIp"`
	HostPort int    `json:"hostPort"`
	EntityID string `json:"entityId"`
	HostID   string `json:"hostId"`
	IdcName  string `json:"idcName"`
	OsType   string `json:"osType,omitempty"`
	Encoding string `json:"encoding,omitempty"` // 系统默认的编码，如果为空，则默认以utf-8值进行处理
	ProxyID  string `json:"proxyId"`
}

//ExecParam 执行参数
type ExecParam struct {
	// 模块名称，支持 script：脚本执行, salt.state：状态应用, file：文件下发
	Pattern string `json:"pattern" validate:"required"`

	// 依据模块名称进行解释
	// pattern为script时，script为脚本内容
	// pattern为salt.state时，script为salt的state内容
	// pattern为file时，script为文件内容或url数组列表
	Script string `json:"script"`
	// 依据pattern进行解释
	// pattern为script时，scriptType为shell, bat, python
	// pattern为file时，scriptType为url或者text
	ScriptType string                 `json:"scriptType" validate:"required"`
	Params     map[string]interface{} `json:"params"`
	RunAs      string                 `json:"runas,omitempty"`
	Password   string                 `json:"password"`
	Timeout    int                    `json:"timeout" validate:"required"`
	Env        map[string]string      `json:"env"`
	ExtendData interface{}            `json:"extendData"`
	// 是否实时输出，像巡检任务、定时任务则不需要实时输出
	RealTimeOutput bool `json:"realTimeOutput"`
}

//Task Exec Result CallBack
type JobCallbackParam struct {
	JobRecordID   string               `json:"jobRecordId"`
	ExecuteID     string               `json:"executeId"`
	ExecuteStatus string               `json:"executeStatus"`
	ResultStatus  string               `json:"resultStatus"`
	HostResults   []HostResultCallback `json:"hostResults"`
}

type HostResultCallback struct {
	EntityID string `json:"entityId"`
	HostIP   string `json:"hostIp"`
	IdcName  string `json:"idcName"`
	Status   string `json:"status"`
	Message  string `json:"message"`
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	Time     string `json:"time"`
}
