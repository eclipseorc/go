package util

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"hash/crc32"
)

// 生成md5
func MD5(b []byte) string {
	c := md5.New()
	c.Write(b)
	return hex.EncodeToString(c.Sum(nil))
}
func MD5s(str string) string { return MD5([]byte(str)) }

//生成sha1
func SHA1(b []byte) string {
	c := sha1.New()
	c.Write([]byte(b))
	return hex.EncodeToString(c.Sum(nil))
}
func SHA1s(str string) string { return SHA1([]byte(str)) }

//生成CRC32
func CRC32(b []byte) uint32 {
	return crc32.ChecksumIEEE(b)
}
func CRC32s(str string) uint32 { return CRC32([]byte(str)) }
