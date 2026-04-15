package customtype

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/zhuofanxu/axb/timex"
)

// JsonTime 定义自定义的时间类型，用于 JSON 序列化和反序列化, 以及数据库读写
type JsonTime time.Time

//goland:noinspection GoMixedReceiverTypes
func (t *JsonTime) MarshalJSON() ([]byte, error) {
	strTime := timex.FormatFromTime(time.Time(*t), timex.TimeTemplateOne)
	return []byte(fmt.Sprintf("\"%s\"", strTime)), nil
}

//goland:noinspection GoMixedReceiverTypes
func (t *JsonTime) UnmarshalJSON(data []byte) error {
	str := string(data)

	// remove the surrounding quotes
	str = strings.Trim(string(data), "\"")

	// 如果是日期格式（长度为10），补充时间部分
	if len(str) == 10 {
		str += " 00:00:00"
	}

	parsedTime, err := time.ParseInLocation(timex.TimeTemplateOne, str, time.Local)
	if err != nil {
		return err
	}
	*t = JsonTime(parsedTime)
	return nil
}

// Value 实现 driver.Valuer 接口
//
//goland:noinspection GoMixedReceiverTypes
func (t JsonTime) Value() (driver.Value, error) {
	return time.Time(t), nil
}

// Scan 实现 sql.Scanner 接口
//
//goland:noinspection GoMixedReceiverTypes
func (t *JsonTime) Scan(value interface{}) error {
	if value == nil {
		*t = JsonTime(time.Time{})
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*t = JsonTime(v)
	case []byte, string:
		var str string
		if b, ok := v.([]byte); ok {
			str = string(b)
		} else {
			str = v.(string)
		}

		tm, err := time.ParseInLocation(timex.TimeTemplateOne, str, time.Local)
		if err != nil {
			return err
		}
		*t = JsonTime(tm)
	default:
		return fmt.Errorf("不支持的类型转换: %T", value)
	}
	return nil
}

// SnowflakeID 定义自定义的 Snowflake ID 类型，用于 JSON 序列化和反序列化
type SnowflakeID uint64

// MarshalJSON implements the json.Marshaler interface
func (id *SnowflakeID) MarshalJSON() ([]byte, error) {
	// 转成字符串，避免 JS 精度丢失
	return []byte(`"` + strconv.FormatUint(uint64(*id), 10) + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (id *SnowflakeID) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)

	// 兼容空值/零值输入
	if s == "" || s == "0" {
		*id = 0
		return nil
	}

	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid id format: %s", s)
	}
	*id = SnowflakeID(v)
	return nil
}

func (id *SnowflakeID) String() string {
	return strconv.FormatUint(uint64(*id), 10)
}
