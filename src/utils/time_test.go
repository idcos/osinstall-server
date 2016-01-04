package utils

import (
	"encoding/json"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestJsonTime(t *testing.T) {
	type person struct {
		Birthday ISOTime `json:"birthday"`
	}
	now, _ := time.Parse(longTimeLayout, "2015-09-24 14:59:06")

	sJSON := `{"birthday":"2015-09-24 14:59:06"}`

	Convey("对象转换成JSON字符串", t, func() {
		p := person{
			Birthday: ISOTime(now),
		}
		bJSON, err := json.Marshal(p)
		So(err, ShouldBeNil)
		So(string(bJSON), ShouldEqual, sJSON)
		So(p.Birthday.String(), ShouldEqual, "2015-09-24 14:59:06")
	})

	Convey("JSON字符串转换成对象", t, func() {
		var p person
		err := json.Unmarshal([]byte(sJSON), &p)
		So(err, ShouldBeNil)
		So(now.Format(longTimeLayout), ShouldEqual, time.Time(p.Birthday).Format(longTimeLayout))
	})

}

func TestUnixSecToISOTime(t *testing.T) {
	Convey("测试unix秒转ISOTime", t, func() {
		jt := UnixSecToISOTime(10) // 10 秒
		So(jt.String(), ShouldEqual, "1970-01-01 08:00:10")
		jt = UnixSecToISOTime(600) // 10 分钟
		So(jt.String(), ShouldEqual, "1970-01-01 08:10:00")
		jt = UnixSecToISOTime(36000) // 10 小时
		So(jt.String(), ShouldEqual, "1970-01-01 18:00:00")
	})
}
