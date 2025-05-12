package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"food-story/pkg/common"
	"food-story/pkg/exceptions"
	"github.com/shopspring/decimal"
	"io"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func IndexToFieldName(indexName string, tableName string) string {
	parts := strings.Split(indexName, "_")

	// Set of parts to ignore
	ignore := map[string]bool{
		"idx":     true,
		"key":     true,
		"pkey":    true,
		"unique":  true,
		tableName: true, // optionally ignore table name
	}

	// Filter out ignored parts
	var filtered []string
	for _, part := range parts {
		if !ignore[part] {
			filtered = append(filtered, part)
		}
	}

	// Convert to camelCase
	for i := range filtered {
		if i == 0 {
			filtered[i] = strings.ToLower(filtered[i]) // first word lowercase
		} else {
			filtered[i] = capitalizeFirst(filtered[i]) // capitalize others
		}
	}

	// Join to camelCase
	return strings.Join(filtered, "")
}

func capitalizeFirst(s string) string {
	for i, r := range s {
		return string(unicode.ToUpper(r)) + s[i+1:]
	}
	return ""
}

func CalculatePageSizeAndNumber(pageSize, pageNumber int64) (int64, int64) {

	if pageSize <= 0 {
		pageSize = int64(common.DefaultPageSize)
	}

	if pageSize > common.MaxPageSize {
		pageSize = int64(common.MaxPageSize)
	}

	if pageNumber <= 0 {
		pageNumber = int64(1)
	}

	return pageSize, (pageNumber - 1) * pageSize
}

func FilterOutZero(slice []int64) []int64 {
	var result []int64
	for _, v := range slice {
		if v != 0 {
			result = append(result, v)
		}
	}
	return result
}

func FilterOutEmptyStr(slice []string) []string {
	var result []string
	for _, v := range slice {
		if v != "" {
			result = append(result, v)
		}
	}
	return result
}

func StrToInt64(value string) (int64, error) {

	if value == "" {
		return 0, exceptions.ErrValueIsEmpty
	}

	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, exceptions.ErrIDInvalidFormat
	}

	return id, nil
}

type SessionData struct {
	SessionID string    `json:"session_id"`
	Expiry    time.Time `json:"expiry"`
}

func EncryptSession(data SessionData, key []byte) (string, error) {
	plaintext, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func DecryptSession(encrypted string, key []byte) (SessionData, error) {
	var data SessionData

	ciphertext, err := base64.URLEncoding.DecodeString(encrypted)
	if err != nil {
		return data, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return data, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return data, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return data, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(plaintext, &data)
	return data, err
}

func ConvertFloatToIntExp(floatNumber float64) int64 {
	num := decimal.NewFromFloat(floatNumber)
	formattedNum := num.StringFixed(2)
	result, _ := strconv.ParseInt(strings.Replace(formattedNum, ".", "", -1), 10, 64)
	return result
}
