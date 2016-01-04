package config

// Loader 定义统一的配置加载接口
type Loader interface {
	Load() (*Config, error)
	Save(*Config) error
}

// Config config 数据结构体
type Config struct {
	Logger struct {
		Color   bool   `ini:"color"`
		Level   string `ini:"level"`
		LogFile string `ini:"logFile"`
	}
	Repo struct {
		Connection string `ini:"connection"`
		Addr       string
	}
	OsInstall struct {
		PxeConfigDir string `ini:"pxeConfigDir"`
	}
}
