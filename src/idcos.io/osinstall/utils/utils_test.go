package utils

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPing(t *testing.T) {
	Convey("www.baidu.com", t, func() {
		So(PingLoop("www.baidu.com", 2, 2), ShouldBeTrue)
	})

	Convey("osinstall.", t, func() {
		So(PingLoop("osinstall.", 2, 2), ShouldBeTrue)
	})

	Convey("www.google.com", t, func() {
		So(PingLoop("www.google.com", 2, 2), ShouldBeFalse)
	})
}
