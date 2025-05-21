package utilities

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	CryptoRand "crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/alexedwards/argon2id"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/sha3"
	"io"
	"math/big"
	"strings"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomString(length int) (string, error) {
	result := make([]byte, length)
	charsetLength := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}

	return string(result), nil
}

func Hash(data string) (string, error) {
	hash, err := argon2id.CreateHash(data, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func VerifyHash(data string, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(data, hash)
	if err != nil {
		return false, err
	}

	return match, nil
}

func EncodeBase64(input []byte) string {
	// Convert string to byte slice (since Base64 encoding works with byte data)
	// Base64 encode the byte slice
	return base64.StdEncoding.EncodeToString(input)
}

func DecodeBase64(text string) ([]byte, error) {
	// Base64 encode the byte slice
	return base64.StdEncoding.DecodeString(text)
}

func EncodeBase64URL(input []byte) string {
	// Example input byte slice
	// Base64 URL encode the byte slice
	base64Url := base64.URLEncoding.EncodeToString(input)

	// Remove padding ('='), as Base64 URL-safe encoding typically omits it
	base64Url = strings.TrimRight(base64Url, "=")

	return base64Url
}

func DecodeBase64URL(text string) ([]byte, error) {
	// Add padding if missing, Base64 URL-safe encoding doesn't require it, but we need to add padding for decoding
	padding := len(text) % 4
	if padding > 0 {
		// Add the necessary padding
		text += strings.Repeat("=", 4-padding)
	}

	return base64.URLEncoding.DecodeString(text)
}

func GenerateAESKey(passphrase []byte, salt []byte) ([]byte, error) {
	// Derive a key using PBKDF2 with SHA3-512
	return pbkdf2.Key(passphrase, salt, 4096, 32, sha3.New512), nil
}

func EncryptAES(key, input []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(input))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(CryptoRand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], input)

	return ciphertext, nil
}

func DecryptAES(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	plaintext := make([]byte, len(ciphertext))
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(plaintext, ciphertext)

	return plaintext, nil
}
