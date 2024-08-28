package timeutils

import (
	"time"
)

const (
	TimeTemplateOne   = "2006-01-02 15:04:05"
	TimeTemplateTwo   = "2006/01/02 15:04:05"
	TimeTemplateThree = "2006-01-02"
	TimeTemplateFour  = "15:04:05"
	TimeTemplateFive  = "2006-01-02 15:04:05.000"
	DaySeconds        = 86400
)

func ParseTimeFromString(value string, tpl string) *time.Time {
	t, err := time.ParseInLocation(tpl, value, time.Local)
	if err != nil {
		return nil
	}
	return &t
}

func ParseTimeFromStringWithLocal(value string, tpl string, local string) *time.Time {
	loc, err := time.LoadLocation(local)
	if err != nil {
		return nil
	}
	t, err := time.ParseInLocation(tpl, value, loc)
	if err != nil {
		return nil
	}
	return &t
}

func FormatFromTime(t time.Time, tpl string) string {
	return t.Format(tpl)
}

func FormatFromTimestamp(timestamp int64, tpl string) string {
	return time.Unix(timestamp, 0).Format(tpl)
}
