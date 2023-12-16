package style

import (
	"time"

	"github.com/aicirt2012/fileintegrity/src/store"
	"github.com/aicirt2012/fileintegrity/src/store/check/style/hierarchy"
	"github.com/aicirt2012/fileintegrity/src/store/check/style/length"
	"github.com/aicirt2012/fileintegrity/src/store/check/style/naming"
	"github.com/aicirt2012/fileintegrity/src/store/file"
	"github.com/aicirt2012/fileintegrity/src/store/ilog"
)

func Check(basePath string, options store.Options) {
	start := time.Now()
	summary := ilog.StyleSummary{}
	logBuffer := ilog.NewAutomaticLogBuffer(basePath, ilog.Style, 10000, options.Log)

	fileHashes := file.LoadContent(basePath)
	summary.HierarchyIssues = hierarchy.Check(fileHashes, &logBuffer)
	summary.NamingIssues = naming.Check(fileHashes, &logBuffer)
	summary.LengthIssues = length.Check(fileHashes, &logBuffer)

	summary.ExecutionTime = time.Since(start)
	logBuffer.Append(summary).Flush()
}
