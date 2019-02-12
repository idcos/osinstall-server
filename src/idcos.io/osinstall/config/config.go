package config

// Loader 定义统一的配置加载接口
type Loader interface {
	Load() (*Config, error)
	Save(*Config) error
}

// Config config 数据结构体
type Config struct {
	Logger struct {
		Color     bool   `ini:"color"`
		Level     string `ini:"level"`
		LogFile   string `ini:"logFile"`
		Logger    Logger
		Repo      Repo
		OsInstall OsInstall
		Vm        Vm
		Rsa       Rsa
		Cron      Cron
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

type Logger struct {
	Color   bool   `ini:"color"`
	Level   string `ini:"level"`
	LogFile string `ini:"logFile"`
}

type Repo struct {
	Connection          string `ini:"connection"`
	ConnectionIsCrypted string `ini:"connectionIsCrypted"`
	Addr                string
}

type OsInstall struct {
	HTTPPort     int
	PxeConfigDir string `ini:"pxeConfigDir"`
}

type Vm struct {
	Storage string `ini:"storage"`
}

type Rsa struct {
	PublicKey  string `ini:"publicKey"`
	PrivateKey string `ini:"privateKey"`
}

type Cron struct {
	InstallTimeout int `ini:"installTimeout"`
}

type ActiveMQ struct {
	Server string `ini:"server"`
}
