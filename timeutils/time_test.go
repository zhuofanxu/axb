package timeutils

import (
	"testing"
	"time"
)

func TestFormatFromTime(t *testing.T) {
	type args struct {
		t   time.Time
		tpl string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "case1", args: args{t: time.Now(), tpl: TimeTemplateThree}, want: "2024-08-28"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatFromTime(tt.args.t, tt.args.tpl); got != tt.want {
				t.Errorf("FormatFromTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
