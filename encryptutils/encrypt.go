package encryptutils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"

	"github.com/zhuofanxu/axb/stringutils"
)

// Md5WithSalt get md5 hash of plainStr, support salts
func Md5WithSalt(plainStr string, salts ...string) string {
	if salts == nil {
		return md5Hash(plainStr)
	} else {
		salt := stringutils.ConcatWithPreallocate(salts...)
		return md5Hash(md5Hash(plainStr+salt) + salt)
	}
}

// Sha256WithSalt get sh256 hash of plainStr, support salts
func Sha256WithSalt(plainStr string, salts ...string) string {
	if salts == nil {
		return sha256Hash(plainStr)
	} else {
		salt := stringutils.ConcatWithPreallocate(salts...)
		return sha256Hash(sha256Hash(plainStr+salt) + salt)
	}
}

func sha256Hash(str string) string {
	h := sha256.New()
	h.Write(stringutils.ConvertToBytes(str))
	return hex.EncodeToString(h.Sum(nil))
}

func md5Hash(str string) string {
	h := md5.New()
	h.Write(stringutils.ConvertToBytes(str))
	return hex.EncodeToString(h.Sum(nil))
}
