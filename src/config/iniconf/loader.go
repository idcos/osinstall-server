package iniconf

import (
	"config"

	"gopkg.in/ini.v1"
)

// INILoader 从ini 文件中加配置数据
type INILoader struct {
	path    string
	content []byte
}

// New 新建INILoader
func New(path string) *INILoader {
	return &INILoader{
		path:    path,
		content: nil,
	}
}

// NewContent 从字符串生成 ini struct for test
func NewContent(content []byte) *INILoader {
	return &INILoader{
		path:    "",
		content: content,
	}
}

// Load 实现Loader 接口 Load()
func (loader *INILoader) Load() (*config.Config, error) {

	var (
		cfg *ini.File
		err error
	)

	if loader.content != nil {
		cfg, err = ini.Load(loader.content)
	} else {
		cfg, err = ini.Load(loader.path)
	}

	if err != nil {
		return nil, err
	}

	var conf = new(config.Config)
	return conf, cfg.MapTo(conf)
}

// Save 实现Loader 接口 Save()
func (loader *INILoader) Save(conf *config.Config) error {
	return nil
}
