package jsonconf

import "testing"

func TestJSONLoader(t *testing.T) {
	loader := New("config_test.json")
	conf, err := loader.Load()
	if err != nil {
		t.Errorf("Load error: %s\n", err)
		return
	}
	if conf.Server.Addr != ":8080" {
		t.Errorf("Config data error\n")
		return
	}
}
