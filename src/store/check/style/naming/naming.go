package naming

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aicirt2012/fileintegrity/src/store/check/style/common"
	"github.com/aicirt2012/fileintegrity/src/store/file"
	"github.com/aicirt2012/fileintegrity/src/store/ilog"
)

var re = regexp.MustCompile(strings.Join([]string{
	` - Copy \(\d{1,3}\)$`,
	` - Kopie \(\d{1,3}\)$`,
	` \(\d{1,3}\)$`,
}, `|`))

func Check(fhs file.FileHashs, logBuffer *ilog.LogFileBuffer) int {
	m := make(common.LogStyleMap)
	for _, fh := range fhs {
		m.PutAll(findPathIssues(fh.RelativePath))
	}
	return m.WriteTo(logBuffer)
}

func findPathIssues(path string) (logs []ilog.StyleLog) {
	sections := common.SplitPath(path)
	for i, section := range sections {
		matches := findSectionIssues(section, i == len(sections)-1)
		if len(matches) > 0 {
			logs = append(logs, ilog.StyleLog{
				IssueType:    ilog.NAMING_ISSUE,
				Reason:       fmt.Sprintf("Path contains copy or rename postfix '%v'", matches[0]),
				RelativePath: common.MinimalPath(i+1, sections),
			})
		}
	}
	if hasExtRepetition(path) {
		logs = append(logs, ilog.StyleLog{
			IssueType:    ilog.NAMING_ISSUE,
			Reason:       "File name contains repeated extension",
			RelativePath: path,
		})
	}
	if hasSpaceSuffix(path) {
		logs = append(logs, ilog.StyleLog{
			IssueType:    ilog.NAMING_ISSUE,
			Reason:       "File name contains space postfix",
			RelativePath: path,
		})
	}
	return logs
}

func findSectionIssues(section string, isFile bool) []string {
	if isFile {
		section = fileNameWithoutExt(section)
	}
	i := re.FindAllString(section, 100)
	if i == nil {
		return []string{}
	}
	return i
}

func hasExtRepetition(path string) bool {
	ext := filepath.Ext(path)
	return len(ext) != 0 && strings.HasSuffix(path, ext+ext)
}

func hasSpaceSuffix(path string) bool {
	ext := filepath.Ext(path)
	return len(ext) != 0 && strings.HasSuffix(path, " "+ext)
}

func fileNameWithoutExt(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
