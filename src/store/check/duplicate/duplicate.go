package duplicate

import (
	"strings"
	"time"

	"github.com/aicirt2012/fileintegrity/src/store"
	"github.com/aicirt2012/fileintegrity/src/store/file"
	"github.com/aicirt2012/fileintegrity/src/store/ilog"
)

func Check(basePath string, options store.Options) {
	start := time.Now()
	summary := ilog.DuplicateSummary{}
	logBuffer := ilog.NewAutomaticLogBuffer(basePath, ilog.Duplicate, 10000, options.Log)

	m, tf, tb := CalcHashSizeMap(file.LoadContent(basePath))
	summary.TotalFiles = tf
	summary.TotalBytes = tb

	df, db, _ := Analyze(m, &logBuffer)
	summary.DuplicateFiles = df
	summary.DuplicateBytes = db

	summary.ExecutionTime = time.Since(start)
	logBuffer.Append(summary).Flush()
}

func Analyze(m UniqueMap, logBuffer *ilog.LogFileBuffer) (int64, int64, []string) {
	var relPaths []string
	var files, bytes int64
	for _, fileHashes := range m.orderedValues() {
		if len(fileHashes) <= 1 {
			continue
		}
		log := ilog.DuplicateLog{
			Hash: fileHashes[0].Hash,
		}
		for i, fileHash := range fileHashes {
			log.AddRelativePath(fileHash.RelativePath)
			if i > 0 {
				files++
				bytes += fileHash.Size
				relPaths = append(relPaths, fileHash.RelativePath)
			}
		}
		logBuffer.Append(log)
	}
	return files, bytes, relPaths
}

// Create map with hash and size as key, ignore small files as well as git files
func CalcHashSizeMap(fileHashs []file.FileHash) (UniqueMap, int64, int64) {
	m := UniqueMap{}
	var totalFiles, totalBytes int64
	for _, fh := range fileHashs {
		totalFiles++
		totalBytes += fh.Size
		if fh.Size <= 100 || strings.Contains(fh.RelativePath, ".git") {
			continue
		}
		m.add(fh)
	}
	return m, totalFiles, totalBytes
}
