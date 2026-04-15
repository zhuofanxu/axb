package timex

import (
	"testing"
	"time"
)

func TestParseTimeFromString(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		tpl     string
		wantNil bool
		check   func(t *testing.T, got *time.Time)
	}{
		{
			name:    "valid datetime",
			value:   "2024-01-15 10:30:00",
			tpl:     TimeTemplateOne,
			wantNil: false,
			check: func(t *testing.T, got *time.Time) {
				if got.Year() != 2024 || got.Month() != 1 || got.Day() != 15 ||
					got.Hour() != 10 || got.Minute() != 30 {
					t.Errorf("unexpected parsed time: %v", got)
				}
			},
		},
		{
			name:    "valid date only",
			value:   "2024-08-28",
			tpl:     TimeTemplateThree,
			wantNil: false,
			check: func(t *testing.T, got *time.Time) {
				if got.Format(TimeTemplateThree) != "2024-08-28" {
					t.Errorf("unexpected date: %v", got)
				}
			},
		},
		{
			name:    "valid slash format",
			value:   "2024/08/28 12:00:00",
			tpl:     TimeTemplateTwo,
			wantNil: false,
		},
		{
			name:    "valid time only",
			value:   "15:04:05",
			tpl:     TimeTemplateFour,
			wantNil: false,
		},
		{
			name:    "invalid value",
			value:   "not-a-date",
			tpl:     TimeTemplateOne,
			wantNil: true,
		},
		{
			name:    "empty string",
			value:   "",
			tpl:     TimeTemplateOne,
			wantNil: true,
		},
		{
			name:    "mismatched template",
			value:   "2024-08-28",
			tpl:     TimeTemplateOne, // 需要 datetime，只给了 date
			wantNil: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseTimeFromString(tt.value, tt.tpl)
			if tt.wantNil {
				if got != nil {
					t.Errorf("expected nil, got %v", got)
				}
				return
			}
			if got == nil {
				t.Fatal("expected non-nil, got nil")
			}
			if tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}

func TestParseTimeFromStringWithLocal(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		tpl     string
		local   string
		wantNil bool
	}{
		{"Asia/Shanghai", "2024-01-15 10:30:00", TimeTemplateOne, "Asia/Shanghai", false},
		{"UTC", "2024-01-15 10:30:00", TimeTemplateOne, "UTC", false},
		{"America/New_York", "2024-01-15 10:30:00", TimeTemplateOne, "America/New_York", false},
		{"invalid timezone", "2024-01-15 10:30:00", TimeTemplateOne, "Invalid/Zone", true},
		{"invalid time value", "bad-time", TimeTemplateOne, "Asia/Shanghai", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseTimeFromStringWithLocal(tt.value, tt.tpl, tt.local)
			if tt.wantNil && got != nil {
				t.Errorf("expected nil, got %v", got)
			}
			if !tt.wantNil && got == nil {
				t.Error("expected non-nil, got nil")
			}
		})
	}

	// 同一字符串在不同时区解析，UTC 时间应不同
	t1 := ParseTimeFromStringWithLocal("2024-01-15 10:00:00", TimeTemplateOne, "Asia/Shanghai")
	t2 := ParseTimeFromStringWithLocal("2024-01-15 10:00:00", TimeTemplateOne, "UTC")
	if t1 == nil || t2 == nil {
		t.Fatal("unexpected nil")
	}
	if t1.Equal(*t2) {
		t.Error("same string in different timezones should yield different UTC instants")
	}
}

func TestFormatFromTime(t *testing.T) {
	fixed := time.Date(2024, 8, 28, 12, 30, 45, 0, time.Local)
	tests := []struct {
		name string
		tpl  string
		want string
	}{
		{"datetime", TimeTemplateOne, "2024-08-28 12:30:45"},
		{"slash datetime", TimeTemplateTwo, "2024/08/28 12:30:45"},
		{"date only", TimeTemplateThree, "2024-08-28"},
		{"time only", TimeTemplateFour, "12:30:45"},
		{"compact", TimeTemplateSix, "20240828123045"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatFromTime(fixed, tt.tpl)
			if got != tt.want {
				t.Errorf("FormatFromTime() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatFromTimestamp(t *testing.T) {
	// 构造本地时间的 Unix 时间戳，再用 FormatFromTimestamp 格式化，结果应一致
	fixed := time.Date(2024, 8, 28, 0, 0, 0, 0, time.Local)
	ts := fixed.Unix()

	got := FormatFromTimestamp(ts, TimeTemplateThree)
	want := "2024-08-28"
	if got != want {
		t.Errorf("FormatFromTimestamp() = %q, want %q", got, want)
	}
}

func TestFormatFromTimestamp_Roundtrip(t *testing.T) {
	// 格式化后再解析，时间应一致（精确到秒）
	original := time.Date(2024, 6, 15, 9, 0, 0, 0, time.Local)
	ts := original.Unix()

	formatted := FormatFromTimestamp(ts, TimeTemplateOne)
	parsed := ParseTimeFromString(formatted, TimeTemplateOne)
	if parsed == nil {
		t.Fatal("roundtrip parse returned nil")
	}
	if !parsed.Equal(original) {
		t.Errorf("roundtrip mismatch: got %v, want %v", parsed, original)
	}
}
