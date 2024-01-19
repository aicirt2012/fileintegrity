package extension

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/aicirt2012/fileintegrity/src/store"
	"github.com/aicirt2012/fileintegrity/src/store/file"
	"github.com/aicirt2012/fileintegrity/src/store/ilog"
)

func Check(basePath string, options store.Options) {
	start := time.Now()
	summary := ilog.ExtensionStatsSummary{}
	logBuffer := ilog.NewAutomaticLogBuffer(basePath, ilog.ExtensionStats, 10000, options.Log)

	m, tb := calcExtensionMap(file.LoadContent(basePath))
	analyze(m, tb, &logBuffer)

	summary.TotalBytes = tb
	summary.UniqueExtensions = len(m)
	summary.ExecutionTime = time.Since(start)
	logBuffer.Append(summary).Flush()
}

func analyze(m extMap, totalBytes int64, logBuffer *ilog.LogFileBuffer) {
	values := m.orderedValues()
	for _, value := range values {
		if value.bytes > 0 {
			value.percentage = float64(value.bytes) / float64(totalBytes) * 100
		}
		logBuffer.Append(ilog.ExtensionStatsLog{
			Name:       value.name,
			Bytes:      value.bytes,
			Percentage: value.percentage,
		})
	}
}

func calcExtensionMap(fileHashs []file.FileHash) (extMap, int64) {
	m := extMap{}
	var totalBytes int64
	for _, fh := range fileHashs {
		totalBytes += fh.Size
		name := strings.ToLower(filepath.Ext(fh.RelativePath))
		if entry, ok := m[name]; ok {
			entry.bytes += fh.Size
		} else {
			m[name] = &ext{
				name:  name,
				bytes: fh.Size,
			}
		}
	}
	return m, totalBytes
}
