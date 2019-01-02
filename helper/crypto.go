package helper

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"hash/crc32"
)

// Md5
func Md5(str string) string {
	hash := md5.New()
	_, _ = hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

// Sha1
func Sha1(str string) string {
	hash := sha1.New()
	_, _ = hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

// Crc32
func Crc32(str string) uint32 {
	return crc32.ChecksumIEEE([]byte(str))
}

// Base64Encode
func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// Base64Decode
func Base64Decode(str string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
