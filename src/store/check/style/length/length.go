package length

import (
	"fmt"

	"github.com/aicirt2012/fileintegrity/src/store/check/style/common"
	"github.com/aicirt2012/fileintegrity/src/store/file"
	"github.com/aicirt2012/fileintegrity/src/store/ilog"
)

const maxPathLen = 260
const maxDirLen = 60

func Check(fhs file.FileHashs, logBuffer *ilog.LogFileBuffer) int {
	m := make(common.LogStyleMap)
	for _, fh := range fhs {
		m.PutAll(findPathIssues(fh.RelativePath))
	}
	return m.WriteTo(logBuffer)
}

func findPathIssues(path string) []ilog.StyleLog {
	logs := []ilog.StyleLog{}
	dirNames := common.SplitDirs(path)
	for i, dirName := range dirNames {
		if len(dirName) > maxDirLen {
			logs = append(logs, ilog.StyleLog{
				IssueType:    ilog.LENGTH_ISSUE,
				Reason:       fmt.Sprintf("Maximum directory length of %v characters exceeded by %v characters", maxDirLen, len(dirName)-maxDirLen),
				RelativePath: common.MinimalPath(i+1, dirNames),
			})
		}
	}
	if len(path) > maxPathLen {
		logs = append(logs, ilog.StyleLog{
			IssueType:    ilog.LENGTH_ISSUE,
			Reason:       fmt.Sprintf("Maximum path length of %v characters exceeded by %v characters", maxPathLen, len(path)-maxPathLen),
			RelativePath: path,
		})
	}
	return logs
}
