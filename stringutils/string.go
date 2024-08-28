package stringutils

import (
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
