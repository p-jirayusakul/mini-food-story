package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"food-story/pkg/common"
	"food-story/pkg/exceptions"
	"io"
	"log/slog"
	"math"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
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

func FilterOutZeroInt(slice []int32) []int32 {
	var result []int32
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

func DecryptSessionToUUID(encrypted string, key []byte) (result uuid.UUID, err error) {
	decrypt, err := DecryptSession(encrypted, key)
	if err != nil {
		slog.Error("DecryptSessionToUUID", "err", err)
		return result, errors.New("invalid session id, please login again")
	}

	if decrypt.SessionID == "" {
		return result, errors.New("session id is empty")
	}

	result, err = uuid.Parse(decrypt.SessionID)
	if err != nil {
		slog.Error("DecryptSessionToUUID", "err", err)
		return result, errors.New("invalid session id, please login again")
	}

	return
}

func ConvertFloatToIntExp(floatNumber float64) int64 {
	num := decimal.NewFromFloat(floatNumber)
	formattedNum := num.StringFixed(2)
	result, _ := strconv.ParseInt(strings.ReplaceAll(formattedNum, ".", ""), 10, 64)
	return result
}

func PgNumericToFloat64(floatNumber pgtype.Numeric) float64 {
	if !floatNumber.Valid {
		return 0
	}

	floatValue, _ := floatNumber.Float64Value()
	return floatValue.Float64
}

func PgTextToStringPtr(text pgtype.Text) *string {
	if !text.Valid {
		return nil
	}

	return &text.String
}

func Float64ToPgNumeric(value float64) pgtype.Numeric {
	return pgtype.Numeric{
		Int:   big.NewInt(ConvertFloatToIntExp(value)),
		Exp:   -2,
		Valid: true,
	}
}

func StringPtrToPgText(value *string) pgtype.Text {
	if value == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{
		String: *value,
		Valid:  true,
	}
}

func UUIDToPgUUID(value uuid.UUID) pgtype.UUID {
	if value == uuid.Nil {
		return pgtype.UUID{Valid: false}
	}
	var byteArray [16]byte = value
	return pgtype.UUID{
		Bytes: byteArray,
		Valid: true,
	}
}

func PgTimestampToThaiISO8601(ts pgtype.Timestamptz) (string, error) {
	if !ts.Valid {
		return "", fmt.Errorf("timestamp is null")
	}

	t := ts.Time

	timeZone := "Asia/Bangkok"
	if os.Getenv("TZ") != "" {
		timeZone = os.Getenv("TZ")
	}

	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		return "", err
	}

	return t.In(loc).Format(time.RFC3339), nil
}

func PareStringToUUID(str string) (uuid.UUID, error) {
	if str == "" {
		return uuid.UUID{}, errors.New("string is empty")
	}
	return uuid.Parse(str)
}

func CalculateTotalPages(totalItems int64, pageSize int64) int64 {
	if pageSize <= 0 {
		pageSize = common.DefaultPageSize
	}

	return int64(math.Ceil(float64(totalItems) / float64(pageSize)))
}

func IsValidTimeZone(timeZone string) bool {
	if timeZone == "" {
		slog.Error("timeZone is empty")
		return false
	}
	_, err := time.LoadLocation(timeZone)
	return err == nil
}
