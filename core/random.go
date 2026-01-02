package core

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"
)

const (
	Letters          = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LettersAndDigits = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	Digits           = "0123456789"
	Alphanumeric     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	HexCharacters    = "0123456789abcdef"
)

func randSequence(n int, charset string) (string, error) {
	if n <= 0 {
		return "", nil
	}

	result := make([]byte, n)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}

	return string(result), nil
}

func GenSalt(length int) (string, error) {
	if length <= 0 {
		length = 50
	}
	return randSequence(length, LettersAndDigits)
}

func GenRandomString(length int, charset string) (string, error) {
	if charset == "" {
		charset = LettersAndDigits
	}
	return randSequence(length, charset)
}

func GenSecureToken(length int) (string, error) {
	if length <= 0 {
		length = 32
	}

	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}

func GenNumericCode(length int) (string, error) {
	if length <= 0 {
		length = 6
	}
	return randSequence(length, Digits)
}

func GenAlphanumeric(length int) (string, error) {
	if length <= 0 {
		length = 16
	}
	return randSequence(length, Alphanumeric)
}

func GenHex(length int) (string, error) {
	if length <= 0 {
		length = 16
	}
	return randSequence(length, HexCharacters)
}

func MustGenSalt(length int) string {
	salt, err := GenSalt(length)
	if err != nil {
		panic("failed to generate salt: " + err.Error())
	}
	return salt
}

func MustGenSecureToken(length int) string {
	token, err := GenSecureToken(length)
	if err != nil {
		panic("failed to generate secure token: " + err.Error())
	}
	return token
}

func MustGenNumericCode(length int) string {
	code, err := GenNumericCode(length)
	if err != nil {
		panic("failed to generate numeric code: " + err.Error())
	}
	return code
}
