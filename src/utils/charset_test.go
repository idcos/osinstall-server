package utils

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCharset(t *testing.T) {
	Convey("utf8字符串与gbk字符串相互转换", t, func() {
		str := "hello, 世界！"         // utf8 string
		gbkStr := UTF82GBK(str)     // utf8-->gbk
		utf8Str := GBK2UTF8(gbkStr) // gbk-->utf8
		So(utf8Str, ShouldEqual, str)
	})

	Convey("将utf8字符串再次转换成utf8字符串", t, func() {
		str := "hello, 世界！"                    // 包含中英文的utf8字符串
		So(GBK2UTF8(str), ShouldNotEqual, str) // 试图将utf8字符串通过'gbk-->utf8'方式转成utf8字符串，不能成立。

		str = "hello, world!"               // 仅仅包含英文的utf8字符串
		So(GBK2UTF8(str), ShouldEqual, str) // 试图将utf8字符串通过'gbk-->utf8'方式转成utf8字符串，成立。
	})
}
