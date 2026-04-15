package customtype

import (
	"encoding/json"
	"testing"
)

func TestSnowflakeID_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    SnowflakeID
		wantErr bool
	}{
		{
			name:    "空字符串应该返回 0",
			input:   `""`,
			want:    0,
			wantErr: false,
		},
		{
			name:    "字符串 0 应该返回 0",
			input:   `"0"`,
			want:    0,
			wantErr: false,
		},
		{
			name:    "有效的 Snowflake ID",
			input:   `"4795027595957637120"`,
			want:    4795027595957637120,
			wantErr: false,
		},
		{
			name:    "合法短数字",
			input:   `"123"`,
			want:    123,
			wantErr: false,
		},
		{
			name:    "合法 - 首位不是 4-9 也应接受",
			input:   `"1234567890123456789"`,
			want:    1234567890123456789,
			wantErr: false,
		},
		{
			name:    "非数字字符串",
			input:   `"abc"`,
			want:    0,
			wantErr: true,
		},
		{
			name:    "超出 uint64 范围",
			input:   `"99999999999999999999"`,
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var id SnowflakeID
			err := json.Unmarshal([]byte(tt.input), &id)

			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && id != tt.want {
				t.Errorf("UnmarshalJSON() got = %v, want %v", id, tt.want)
			}
		})
	}
}
