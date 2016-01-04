package jsonconf

import (
	"config"
	"encoding/json"
	"os"
)

// JSONLoader 从json 文件中加配置数据
type JSONLoader struct {
	path string
}

// New 新建JSONLoader
func New(path string) *JSONLoader {
	return &JSONLoader{path}
}

// Load 实现Loader 接口 Load()
func (loader *JSONLoader) Load() (*config.Config, error) {
	f, err := os.Open(loader.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var conf config.Config

	return &conf, json.NewDecoder(f).Decode(&conf)
}

// Save 实现Loader 接口 Save()
func (loader *JSONLoader) Save(conf *config.Config) error {
	return nil
}
