package crypto

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"

	"github.com/zhuofanxu/axb/stringx"
)

// Md5WithSalt get md5 hash of plainStr, support salts
func Md5WithSalt(plainStr string, salts ...string) string {
	if salts == nil {
		return md5Hash(plainStr)
	} else {
		salt := stringx.ConcatWithPreallocate(salts...)
		return md5Hash(md5Hash(plainStr+salt) + salt)
	}
}

// Sha256WithSalt get sh256 hash of plainStr, support salts
func Sha256WithSalt(plainStr string, salts ...string) string {
	if salts == nil {
		return sha256Hash(plainStr)
	} else {
		salt := stringx.ConcatWithPreallocate(salts...)
		return sha256Hash(sha256Hash(plainStr+salt) + salt)
	}
}

// HashPassword 使用 bcrypt 加密密码
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// CheckPassword 验证密码是否匹配
func CheckPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

func sha256Hash(str string) string {
	h := sha256.New()
	h.Write(stringx.ConvertToBytes(str))
	return hex.EncodeToString(h.Sum(nil))
}

func md5Hash(str string) string {
	h := md5.New()
	h.Write(stringx.ConvertToBytes(str))
	return hex.EncodeToString(h.Sum(nil))
}
