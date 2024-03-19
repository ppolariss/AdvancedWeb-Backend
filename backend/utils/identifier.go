package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func MakePassword(rawPassword string) string {
	return MakeMD5(rawPassword)
}

func CheckPassword(rawPassword, encryptPassword string) bool {
	return MakeMD5(rawPassword) == encryptPassword
}

func MakeMD5(raw string) string {
	hash := md5.New()
	hash.Write([]byte(raw))
	return hex.EncodeToString(hash.Sum(nil))
}
