package stringx

import (
	"strconv"
	"strings"
	"unsafe"
)

// ConvertToBytes convert string to bytes efficiently without memory allocation
func ConvertToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

// ConvertToString convert bytes to string efficiently without memory allocation
func ConvertToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// ConcatWithPreallocate concat strings efficiently
func ConcatWithPreallocate(strs ...string) string {
	builder := strings.Builder{}
	preCap := 0
	for _, str := range strs {
		// pre-allocated memory
		preCap += len(str)
	}
	builder.Grow(preCap)
	for _, str := range strs {
		builder.WriteString(str)
	}
	return builder.String()
}

func ConcatWithBuilder(strs ...string) string {
	builder := strings.Builder{}
	for _, str := range strs {
		builder.WriteString(str)
	}
	return builder.String()
}

func ConcatWithPlus(strs ...string) string {
	var s string
	for _, str := range strs {
		s += str
	}
	return s
}

func ParseHexToUint64(s string) (uint, error) {
	v := strings.TrimSpace(s)
	if strings.HasPrefix(v, "0x") || strings.HasPrefix(v, "0X") {
		v = v[2:]
	}
	// 空字符串或非十六进制将报错
	u64, err := strconv.ParseUint(v, 16, 64)
	if err != nil {
		return 0, err
	}
	return uint(u64), nil
}
