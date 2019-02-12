package utils

import (
	"fmt"
	"strings"
	"time"
)

// ISOTime 标准库Time的别名类型。用于格式化为ISO8601标准的时间（形如'2015-09-09 09:09:09'）。
type ISOTime time.Time

// longTimeLayout 时间格式
var longTimeLayout = "2006-01-02 15:04:05"

// MarshalJSON 序列化，实现json.Marshaller接口
func (t ISOTime) MarshalJSON() ([]byte, error) {
	// 注意：序列化成的时间字符串必须包含在双引号当中
	return []byte(fmt.Sprintf("%q", time.Time(t).Format(longTimeLayout))), nil
}

// UnmarshalJSON 反序列化,实现json.Unmarshaler接口
func (t *ISOTime) UnmarshalJSON(b []byte) error {
	sTime := fmt.Sprintf("%s", string(b))

	if strings.HasPrefix(sTime, "\"") {
		sTime = sTime[1:]
	}

	if strings.HasSuffix(sTime, "\"") {
		sTime = sTime[0 : len(sTime)-1]
	}

	tmpT, err := time.Parse(longTimeLayout, sTime)
	if err != nil {
		return err
	}

	*t = ISOTime(tmpT)
	return nil
}

// MarshalYAML 序列化成YAML
func (t ISOTime) MarshalYAML() (interface{}, error) {
	return time.Time(t).Format(longTimeLayout), nil
}

// UnmarshalYAML YAML反序列化
// func (t ISOTime) UnmarshalYAML(unmarshal func(interface{}) error) error {
// 	return nil
// }

func (t ISOTime) String() string {
	return time.Time(t).Format(longTimeLayout)
}

// UnixSecToISOTime UNIX秒转化成ISOTime类型
func UnixSecToISOTime(unixSecond int64) ISOTime {
	return ISOTime(time.Unix(unixSecond, 0))
}
