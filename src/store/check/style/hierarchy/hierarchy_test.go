package hierarchy

import (
	"testing"

	"github.com/aicirt2012/fileintegrity/src/store/file"
	"github.com/aicirt2012/fileintegrity/src/store/ilog"
	"github.com/aicirt2012/fileintegrity/tests/common"
	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {
	fhs := file.FileHashs{
		{RelativePath: common.NormalizePath("root/simple path/f.txt")},
		{RelativePath: common.NormalizePath("root/linux ubuntu/another linux/ubuntu/f1.txt")},
		{RelativePath: common.NormalizePath("root/linux ubuntu/another linux/ubuntu/f2.txt")},
	}
	logBuffer := ilog.NewManualLogBuffer("", ilog.Style, ilog.Options{})
	actual := Check(fhs, &logBuffer)
	assert.Equal(t, 1, actual)
}

func TestDuplicateWords(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Too short for duplicate",
			input:    "root/unix/another unix/f.txt",
			expected: []string{},
		},
		{
			name:     "One duplicate words",
			input:    "root/linux/another linux/f.txt",
			expected: []string{"linux"},
		},
		{
			name:     "Tow duplicate words",
			input:    "root/linux ubuntu/another linux/ubuntu/f.txt",
			expected: []string{"linux", "ubuntu"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := duplicateWords(common.NormalizePath(c.input))
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestPathWords(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Convert to lower case words",
			input:    "root/sub x/file Name.txt",
			expected: []string{"root", "sub", "x", "file", "name.txt"},
		},
		{
			name:     "No empty word",
			input:    "root/sub x/",
			expected: []string{"root", "sub", "x"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := pathWords(common.NormalizePath(c.input))
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestMinimalPath(t *testing.T) {
	cases := []struct {
		name       string
		path       string
		duplicates []string
		expected   string
	}{
		{
			name:       "Case sensitive path",
			path:       "/hallo/duplicate/another Duplicate/prune path/any name.txt",
			duplicates: []string{"duplicate"},
			expected:   "/hallo/duplicate/another Duplicate",
		},
		{
			name:       "Case sensitive path with postfix",
			path:       "/hallo/duplicate/another Duplicate postfix/prune path/any name.txt",
			duplicates: []string{"duplicate"},
			expected:   "/hallo/duplicate/another Duplicate postfix",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := minimalPath(common.NormalizePath(c.path), c.duplicates)
			assert.Equal(t, common.NormalizePath(c.expected), actual)
		})
	}
}

func TestILastIndex(t *testing.T) {
	assert.Equal(t, 2, iLastIndex("aaa", "a"))
	assert.Equal(t, 2, iLastIndex("aaa", "A"))
	assert.Equal(t, 2, iLastIndex("AAA", "a"))
	assert.Equal(t, 2, iLastIndex("AAA", "A"))
}

func TestFirstIndexAfter(t *testing.T) {
	cases := []struct {
		name     string
		s        string
		subStr   string
		idx      int
		expected int
	}{
		{
			name:     "Normal case",
			s:        "text with - postfix",
			subStr:   "-",
			idx:      4,
			expected: 10,
		},
		{
			name:     "Edge case",
			s:        " -",
			subStr:   "-",
			idx:      1,
			expected: 1,
		},
		{
			name:     "Out of range",
			s:        " -",
			subStr:   "-",
			idx:      2,
			expected: -1,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := firstIndexAfter(c.s, c.subStr, c.idx)
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestUniqueDirPaths(t *testing.T) {
	input := file.FileHashs{
		{RelativePath: common.NormalizePath("root/sub/f1.txt")},
		{RelativePath: common.NormalizePath("root/sub/f2.txt")},
	}
	expected := []string{common.NormalizePath("root/sub")}
	actual := uniqueDirPaths(input)
	assert.Equal(t, expected, actual)
}

func TestFmtReason(t *testing.T) {
	input := []string{"a", "b", "c"}
	expected := "Path contains the words 'a', 'b' and 'c' on multiple hierarchies"
	actual := fmtReason(input)
	assert.Equal(t, expected, actual)
}
