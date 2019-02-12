package iniconf

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestINILoader(t *testing.T) {
	loader := New("config_test.ini")
	conf, err := loader.Load()
	if err != nil {
		t.Errorf("Load error: %s\n", err)
		return
	}

	Convey("Logger config 校验", t, func() {
		So(conf.Logger.Color, ShouldEqual, true)
	})

	Convey("Repo config 校验", t, func() {
		So(conf.Repo.Connection, ShouldEqual, "admin:admin@10.0.2.8/wyb-devdb?charset=utf8&parseTime=True&loc=Local")
	})
}

func TestINILoaderContent(t *testing.T) {
	var iniContent = `[Logger]
color = true

[Repo]
connection = "admin:admin@10.0.2.8/wyb-devdb?charset=utf8&parseTime=True&loc=Local"`

	loader := NewContent([]byte(iniContent))
	conf, err := loader.Load()
	if err != nil {
		t.Errorf("Load error: %s\n", err)
		return
	}

	Convey("Logger config 校验", t, func() {
		So(conf.Logger.Color, ShouldEqual, true)
	})

	Convey("Repo config 校验", t, func() {
		So(conf.Repo.Connection, ShouldEqual, "admin:admin@10.0.2.8/wyb-devdb?charset=utf8&parseTime=True&loc=Local")
	})
}
