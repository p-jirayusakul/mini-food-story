package utils

import (
	"food-story/pkg/common"
	"strings"
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
