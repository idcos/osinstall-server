package agent

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestExecScript(t *testing.T) {
	Convey("执行命令: ls -l", t, func() {
		_, err := execScript("ls -l")
		So(err, ShouldBeNil)
	})

	Convey("执行命令: ls -l / | grep bin", t, func() {
		_, err := execScript("ls -l / | grep bin")
		So(err, ShouldBeNil)
	})
}
