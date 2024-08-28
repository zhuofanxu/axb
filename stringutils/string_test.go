package stringutils

import (
	"reflect"
	"testing"
)

func TestConvertToBytes(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{name: "case1", args: args{s: "abcd"}, want: []byte{97, 98, 99, 100}},
		{name: "case2", args: args{s: "你好！"}, want: []byte{228, 189, 160, 229, 165, 189, 239, 188, 129}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertToBytes(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
