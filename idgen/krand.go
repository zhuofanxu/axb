package idgen

import (
	"crypto/rand"
	"math/big"
)

const (
	RandKindNum        = 0 // 纯数字
	RandKindLower      = 1 // 小写字母
	RandKindUpper      = 2 // 大写字母
	RandKindNumLower   = 4 // 数字小写混合
	RandKindNumUpper   = 5 // 数字大写混合
	RandKindLowerUpper = 6 // 大小写混合
	RandKindAll        = 7 // 数字字母混合
)

var (
	// 排除易混淆字符：0,1,O,o,l
	digits       = []byte("23456789")
	upperLetters = []byte("ABCDEFGHIJKLMNPQRSTUVWXYZ")
	lowerLetters = []byte("abcdefghijkmnpqrstuvwxyz")
)

// GenKrandId 使用加密安全随机源生成ID
func GenKrandId(size int, kind int) string {
	result := make([]byte, size)

	switch kind {
	case RandKindNum:
		fillRandomChars(result, digits)
	case RandKindLower:
		fillRandomChars(result, lowerLetters)
	case RandKindUpper:
		fillRandomChars(result, upperLetters)
	case RandKindNumLower:
		chars := append(digits, lowerLetters...)
		fillRandomChars(result, chars)
	case RandKindNumUpper:
		chars := append(digits, upperLetters...)
		fillRandomChars(result, chars)
	case RandKindLowerUpper:
		chars := append(upperLetters, lowerLetters...)
		fillRandomChars(result, chars)
	case RandKindAll:
		chars := append(append(digits, upperLetters...), lowerLetters...)
		fillRandomChars(result, chars)
	}

	return string(result)
}

func fillRandomChars(dst []byte, chars []byte) {
	for i := range dst {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		dst[i] = chars[n.Int64()]
	}
}
