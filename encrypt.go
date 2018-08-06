package encdec

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"os"
)

func Encrypt(content []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(content))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], content)

	return ciphertext, nil
}
func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}
func Exists(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	}
	return false
}
func IsDir(file string) bool {
	if stat, err := os.Stat(file); err == nil && stat.IsDir() {
		return true
	}
	return false
}
func EncryptString(input string, key []byte) (out string, err error) {
	ciphertext, err := Encrypt([]byte(input), key)

	if err == nil {
		out = hex.EncodeToString(ciphertext)
	}

	return out, err
}
func DecryptString(input string, key []byte) (out string, err error) {
	ciphertext, _ := hex.DecodeString(input)

	ciphertext, err = Decrypt(ciphertext, key)

	if err == nil {
		out = string(ciphertext)
	}
	return out, err
}
func Md5Sum(content []byte) string {
	return hex.EncodeToString(byte2string(md5.Sum(content)))
}
func byte2string(in [16]byte) []byte {
	return in[:16]
}
