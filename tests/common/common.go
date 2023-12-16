package common

import (
	"path/filepath"
	"strings"
)

// Normalizes static test data to fit the os specific notation
func NormalizePath(path string) string {
	path = strings.ReplaceAll(path, `\`, string(filepath.Separator))
	return strings.ReplaceAll(path, `/`, string(filepath.Separator))
}

func StaticContent(sizeInkB int) string {
	n := sizeInkB * 1024
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteString("A")
	}
	return sb.String()
}
