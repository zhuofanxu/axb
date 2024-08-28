package encryptutils

import (
	"testing"
)

func TestSha256WithSalt(t *testing.T) {
	type args struct {
		plainStr string
		salts    []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "case1", args: args{plainStr: "123456", salts: []string{"salt1", "salt2"}}, want: "2396b938f0096392ba0a32c35577113c427cc68e890a254466b74d1195f0545d"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sha256WithSalt(tt.args.plainStr, tt.args.salts...); got != tt.want {
				t.Errorf("Sha256WithSalt() = %v, want %v", got, tt.want)
			}
		})
	}
}
