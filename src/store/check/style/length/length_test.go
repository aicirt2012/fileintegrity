package length

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aicirt2012/fileintegrity/src/store/file"
	"github.com/aicirt2012/fileintegrity/src/store/ilog"
	"github.com/aicirt2012/fileintegrity/tests/common"
	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {
	fhs := file.FileHashs{
		{RelativePath: extendTo("root/l#ng dir name/f.txt", maxPathLen+1)},
	}
	logBuffer := ilog.NewManualLogBuffer("", ilog.Style, ilog.Options{})
	actual := Check(fhs, &logBuffer)
	assert.Equal(t, 2, actual)
}

func TestFindPathIssues(t *testing.T) {
	cases := []struct {
		name     string
		path     string
		expected []ilog.StyleLog
	}{
		{
			name:     "Max valid section len",
			path:     extendBy("root/#/f.txt", maxDirLen),
			expected: []ilog.StyleLog{},
		},
		{
			name: "Section exceeds max len",
			path: extendBy("root/#/f.txt", maxDirLen+1),
			expected: []ilog.StyleLog{
				{
					IssueType:    ilog.LENGTH_ISSUE,
					Reason:       fmt.Sprintf("Maximum directory length of %v characters exceeded by 1 characters", maxDirLen),
					RelativePath: common.NormalizePath(extendBy("root/#", maxDirLen+1)),
				},
			},
		},
		{
			name:     "Max valid path len",
			path:     extendTo("root/#/#/#/#/#/f.txt", maxPathLen),
			expected: []ilog.StyleLog{},
		},
		{
			name: "Path exceeds max len",
			path: extendTo("root/#/#/#/#/#/f.txt", maxPathLen+1),
			expected: []ilog.StyleLog{
				{
					IssueType:    ilog.LENGTH_ISSUE,
					Reason:       fmt.Sprintf("Maximum path length of %v characters exceeded by 1 characters", maxPathLen),
					RelativePath: common.NormalizePath(extendTo("root/#/#/#/#/#/f.txt", maxPathLen+1)),
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := findPathIssues(common.NormalizePath(c.path))
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestExtendTo(t *testing.T) {
	assert.Equal(t, "looong", extendTo("l#ng", 6))
	assert.Equal(t, "loonog", extendTo("l#n#g", 6))
	assert.Equal(t, "looonoogoo", extendTo("l#n#g#", 10))
}

func extendTo(s string, n int) string {
	p := "#"
	nrP := strings.Count(s, p)
	sLen := len(s) - nrP
	if sLen >= n {
		panic("input is already too long")
	}
	fullRepetitions := (n - sLen) / nrP
	partialRepetitions := (n - sLen) % nrP
	for i := 0; i < nrP; i++ {
		replacement := strings.Repeat("o", fullRepetitions)
		if partialRepetitions > 0 {
			replacement += "o"
			partialRepetitions--
		}
		s = strings.Replace(s, p, replacement, 1)
	}
	return s
}

func TestExtendBy(t *testing.T) {
	assert.Len(t, extendBy("l#ng", 7), 10)
}

func extendBy(s string, n int) string {
	return strings.Replace(s, "#", strings.Repeat("o", n), 1)
}
