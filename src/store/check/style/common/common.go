package common

import (
	"fmt"
	"path/filepath"
	"strings"
)

func SplitPath(path string) []string {
	return strings.Split(path, string(filepath.Separator))
}

func SplitDirs(path string) []string {
	return SplitPath(filepath.Dir(path))
}

func MinimalPath(i int, sections []string) string {
	return strings.Join(sections[:i], string(filepath.Separator))
}

func JoinWithAnd(item []string) string {
	length := len(item)

	if length == 0 {
		return ""
	}
	if length == 1 {
		return item[0]
	}
	return strings.Join(item[:length-1], ", ") + " and " + item[length-1]
}

func Quote(words []string) []string {
	q := []string{}
	for _, word := range words {
		q = append(q, fmt.Sprintf(`'%s'`, word))
	}
	return q
}

func PluralS(arr []string) string {
	if len(arr) > 1 {
		return "s"
	}
	return ""
}
