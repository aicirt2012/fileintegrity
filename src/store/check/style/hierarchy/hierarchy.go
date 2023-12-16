package hierarchy

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/aicirt2012/fileintegrity/src/store/check/style/common"
	"github.com/aicirt2012/fileintegrity/src/store/file"
	"github.com/aicirt2012/fileintegrity/src/store/ilog"
	"golang.org/x/exp/maps"
)

func Check(fhs file.FileHashs, logBuffer *ilog.LogFileBuffer) int {
	m := make(common.LogStyleMap)
	for _, dir := range uniqueDirPaths(fhs) {
		duplicates := duplicateWords(dir)
		if len(duplicates) > 0 {
			m.Put(ilog.StyleLog{
				IssueType:    ilog.HIERARCHY_ISSUE,
				Reason:       fmtReason(duplicates),
				RelativePath: minimalPath(dir, duplicates),
			})
		}
	}
	return m.WriteTo(logBuffer)
}

func duplicateWords(path string) []string {
	duplicates := []string{}
	existing := make(map[string]bool)
	for _, word := range pathWords(path) {
		if len(word) < 5 {
			continue
		}
		if _, ok := existing[word]; ok {
			duplicates = append(duplicates, word)
		}
		existing[word] = true
	}
	return duplicates
}

func pathWords(path string) []string {
	path = strings.ToLower(path)
	path = strings.ReplaceAll(path, string(filepath.Separator), " ")
	words := []string{}
	for _, w := range strings.Split(path, " ") {
		if w != "" {
			words = append(words, w)
		}
	}
	return words
}

func minimalPath(path string, duplicates []string) string {
	index := 0
	for _, duplicate := range duplicates {
		currentIndex := iLastIndex(path, duplicate) + len(duplicate)
		index = max(index, currentIndex)
	}
	index = firstIndexAfter(path, string(filepath.Separator), index)
	if index == -1 {
		return path
	}
	return path[:index]
}

func iLastIndex(s string, subStr string) int {
	s = strings.ToLower(s)
	subStr = strings.ToLower(subStr)
	return strings.LastIndex(s, subStr)
}

func firstIndexAfter(s string, subStr string, idx int) int {
	if len(s) < idx {
		return -1
	}
	s = s[idx:]
	result := strings.Index(s, subStr)
	if result == -1 {
		return -1
	}
	return idx + result
}

func uniqueDirPaths(fhs file.FileHashs) []string {
	m := make(map[string]bool)
	for _, fh := range fhs {
		m[filepath.Dir(fh.RelativePath)] = true
	}
	values := maps.Keys(m)
	return values
}

func fmtReason(duplicates []string) string {
	duplicates = common.Quote(duplicates)
	words := common.JoinWithAnd(duplicates)
	s := common.PluralS(duplicates)
	return fmt.Sprintf("Path contains the word%v %v on multiple hierarchies", s, words)
}
