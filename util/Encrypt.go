package util

import (
	"bytes"
	"crypto/aes"
	"crypto/md5"
	"encoding/hex"
	"errors"
)

// md5编码
func Md5(s string) string {
	b := []byte(s)
	checksum := md5.Sum(b)
	return hex.EncodeToString(checksum[:])
}

// AES加密
// cipherKey为AES加密密钥
func AESEncrypt(cipherKey string, src string, pkcsType int) (error, string) {
	if pkcsType == 0 {
		pkcsType = 7
	}
	cb := []byte(cipherKey)
	block, err := aes.NewCipher(cb)
	if err != nil {
		return err, ""
	}
	bs := block.BlockSize()
	sb := []byte(src)
	if pkcsType == 5 {
		sb = pkcs5Pad(sb, bs)
	} else if pkcsType == 7 {
		sb = pkcs7Pad(sb, bs)
	}
	r := make([]byte, len(sb))
	dst := r
	for len(sb) > 0 {
		block.Encrypt(dst, sb)
		sb = sb[bs:]
		dst = dst[bs:]
	}

	dst = make([]byte, hex.EncodedLen(len(r)))
	hex.Encode(dst, r)
	s := string(dst)
	return nil, s
}

func AESDecrypt(cipherKey string, cipherText string, pkcsType int) (string, error) {
	if pkcsType == 0 {
		pkcsType = 7
	}
	b := []byte(cipherKey)
	block, err := aes.NewCipher(b)
	if err != nil {
		return "", err
	}
	src := make([]byte, hex.DecodedLen(len(cipherText)))
	tb := []byte(cipherText)
	_, err = hex.Decode(src, tb)
	if err != nil {
		return "", err
	}
	bs := block.BlockSize()
	r := make([]byte, len(src))
	dst := r
	for len(src) > 0 {
		block.Decrypt(dst, src)
		src = src[bs:]
		dst = dst[bs:]
	}
	var res []byte
	if pkcsType == 5 {
		res, _ = pkcs5UnPad(r)
	} else if pkcsType == 7 {
		res = pkcs7UnPad(r)
	}
	return string(res), nil
}

// pkcs5填充字节，数据末尾不满16字节的数据，缺多少个字节，就填充多少个字节的几
func pkcs5Pad(d []byte, bs int) []byte {
	padSize := ((len(d) / bs) + 1) * bs
	pad := padSize - len(d)
	for i := len(d); i < padSize; i++ {
		d = append(d, byte(pad))
	}
	return d
}

// pkcs5移除填充字节
func pkcs5UnPad(b []byte) ([]byte, error) {
	l := len(b)
	if l == 0 {
		return nil, errors.New("input []byte is empty")
	}
	last := int(b[l-1])
	pad := b[l-last : l]
	isPad := true
	for _, v := range pad {
		if int(v) != last {
			isPad = false
			break
		}
	}
	if !isPad {
		return b, errors.New("remove pad error")
	}
	return b[:l-last], nil
}

// pkcs7填充字节
func pkcs7Pad(d []byte, bs int) []byte {
	pad := bs - len(d)%bs
	padding := bytes.Repeat([]byte{byte(pad)}, pad)
	return append(d, padding...)
}

func pkcs7UnPad(b []byte) []byte {
	length := len(b)
	padLen := int(b[length-1])
	return b[:(length - padLen)]
}
