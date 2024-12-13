package utils

import (
	"path/filepath"
	"strings"
)

func GetStringPtr(s string) *string {
	return &s
}

func RemoveBackticks(text string) string {
	return strings.Trim(text, "`")
}

func RemoveFileExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
