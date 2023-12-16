package naming

import (
	"testing"

	"github.com/aicirt2012/fileintegrity/src/store/file"
	"github.com/aicirt2012/fileintegrity/src/store/ilog"
	"github.com/aicirt2012/fileintegrity/tests/common"
	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {
	fhs := file.FileHashs{
		{RelativePath: common.NormalizePath("root/s/a .txt")},
		{RelativePath: common.NormalizePath("root/s/a.txt.txt")},
		{RelativePath: common.NormalizePath("root/x (1)/c.txt")},
		{RelativePath: common.NormalizePath("root/x (1)/b.txt")},
	}
	logBuffer := ilog.NewManualLogBuffer("", ilog.Style, ilog.Options{})
	actual := Check(fhs, &logBuffer)
	assert.Equal(t, 3, actual)
}

func TestFindPathIssues(t *testing.T) {
	cases := []struct {
		name     string
		path     string
		expected []ilog.StyleLog
	}{
		{
			name: "Simple case",
			path: "root/sub - Copy (1)/f.txt",
			expected: []ilog.StyleLog{
				{
					IssueType:    ilog.NAMING_ISSUE,
					Reason:       "Path contains copy or rename postfix ' - Copy (1)'",
					RelativePath: common.NormalizePath("root/sub - Copy (1)"),
				},
			},
		},
		{
			name: "Multiple section case",
			path: "root/s (1)/name (2).txt",
			expected: []ilog.StyleLog{
				{
					IssueType:    ilog.NAMING_ISSUE,
					Reason:       "Path contains copy or rename postfix ' (1)'",
					RelativePath: common.NormalizePath("root/s (1)"),
				},
				{
					IssueType:    ilog.NAMING_ISSUE,
					Reason:       "Path contains copy or rename postfix ' (2)'",
					RelativePath: common.NormalizePath("root/s (1)/name (2).txt"),
				},
			},
		},
		{
			name: "Section and file space case",
			path: "root/s (1)/name .txt",
			expected: []ilog.StyleLog{
				{
					IssueType:    ilog.NAMING_ISSUE,
					Reason:       "Path contains copy or rename postfix ' (1)'",
					RelativePath: common.NormalizePath("root/s (1)"),
				},
				{
					IssueType:    ilog.NAMING_ISSUE,
					Reason:       "File name contains space postfix",
					RelativePath: common.NormalizePath("root/s (1)/name .txt"),
				},
			},
		},
		{
			name: "Section and file extension case",
			path: "root/s (1)/name.txt.txt",
			expected: []ilog.StyleLog{
				{
					IssueType:    ilog.NAMING_ISSUE,
					Reason:       "Path contains copy or rename postfix ' (1)'",
					RelativePath: common.NormalizePath("root/s (1)"),
				},
				{
					IssueType:    ilog.NAMING_ISSUE,
					Reason:       "File name contains repeated extension",
					RelativePath: common.NormalizePath("root/s (1)/name.txt.txt"),
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			input := common.NormalizePath(c.path)
			actual := findPathIssues(input)
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestFindSectionIssues(t *testing.T) {
	cases := []struct {
		name     string
		section  string
		isFile   bool
		expected []string
	}{
		{
			name:     "Rename count lower bound",
			section:  "any (1)",
			isFile:   false,
			expected: []string{" (1)"},
		},
		{
			name:     "Rename count upper bound",
			section:  "any (999)",
			isFile:   false,
			expected: []string{" (999)"},
		},
		{
			name:     "Rename count out of bound",
			section:  "any (1000)",
			isFile:   false,
			expected: []string{},
		},
		{
			name:     "English copy lower bound",
			section:  "any - Copy (1)",
			isFile:   false,
			expected: []string{" - Copy (1)"},
		},
		{
			name:     "English copy upper bound",
			section:  "any - Copy (999)",
			isFile:   false,
			expected: []string{" - Copy (999)"},
		},
		{
			name:     "English copy out of bound",
			section:  "any - Copy (1000)",
			isFile:   false,
			expected: []string{},
		},
		{
			name:     "German copy lower bound",
			section:  "any - Kopie (1)",
			isFile:   false,
			expected: []string{" - Kopie (1)"},
		},
		{
			name:     "German copy upper bound",
			section:  "any - Kopie (999)",
			isFile:   false,
			expected: []string{" - Kopie (999)"},
		},
		{
			name:     "German copy out of bound",
			section:  "any - Kopie (1000)",
			isFile:   false,
			expected: []string{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := findSectionIssues(c.section, c.isFile)
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestHasExtRepetition(t *testing.T) {
	assert.Equal(t, false, hasExtRepetition("filename"))
	assert.Equal(t, false, hasExtRepetition("filename.txt"))
	assert.Equal(t, true, hasExtRepetition("filename.txt.txt"))
}

func TestHasSpaceSuffix(t *testing.T) {
	assert.Equal(t, false, hasSpaceSuffix("filename"))
	assert.Equal(t, false, hasSpaceSuffix("filename.txt"))
	assert.Equal(t, true, hasSpaceSuffix("filename .txt"))
}

func TestFileNameWithoutExt(t *testing.T) {
	assert.Equal(t, "filename", fileNameWithoutExt("filename.txt"))
	assert.Equal(t, "filename", fileNameWithoutExt("filename"))
}
