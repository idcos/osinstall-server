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
		Connection          string `ini:"connection"`
		ConnectionIsCrypted string `ini:"connectionIsCrypted"`
		Addr                string
	}
	OsInstall struct {
		PxeConfigDir string `ini:"pxeConfigDir"`
	}
	Vm struct {
		Storage string `ini:"storage"`
	}
	Rsa struct {
		PublicKey  string `ini:"publicKey"`
		PrivateKey string `ini:"privateKey"`
	}
	Cron struct {
		InstallTimeout int `ini:"installTimeout"`
	}
	ActiveMQ struct {
		Server string `ini:"server"`
	}
}
