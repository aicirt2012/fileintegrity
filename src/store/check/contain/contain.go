package contain

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/aicirt2012/fileintegrity/src/store"
	"github.com/aicirt2012/fileintegrity/src/store/check/duplicate"
	"github.com/aicirt2012/fileintegrity/src/store/file"
	"github.com/aicirt2012/fileintegrity/src/store/ilog"
	"golang.org/x/exp/maps"
)

func Check(basePath string, externalPath string, fix bool, options store.Options) {
	start := time.Now()
	summary := ilog.ContainedSummary{}
	logBuffer := ilog.NewAutomaticLogBuffer(basePath, ilog.Contains, 10000, options.Log)

	options.Backup = true
	err := store.Upsert(externalPath, options)
	if err != nil {
		log.Fatal(err)
	}
	baseM, _, _ := duplicate.CalcHashSizeMap(file.LoadContent(basePath))
	externalM, tf, tb := duplicate.CalcHashSizeMap(file.LoadContent(externalPath))
	summary.TotalFiles = tf
	summary.TotalBytes = tb

	cf, cb, containedRelPaths := analyze(baseM, externalM, &logBuffer)
	summary.ContainedFiles = cf
	summary.ContainedBytes = cb

	df, db, duplicateRelPaths := duplicate.Analyze(externalM, &logBuffer)
	summary.DuplicateFiles = df
	summary.DuplicateBytes = db

	if fix {
		removeFiles(externalPath, append(containedRelPaths, duplicateRelPaths...))
	}

	summary.ExecutionTime = time.Since(start)
	logBuffer.Append(summary).Flush()
}

func analyze(baseM duplicate.UniqueMap, externalM duplicate.UniqueMap, logBuffer *ilog.LogFileBuffer) (int64, int64, []string) {
	var relPaths []string
	var files, bytes int64
	for _, externalKey := range maps.Keys(externalM) {
		if !baseM.Has(externalKey) {
			continue
		}
		baseFileHashs := baseM[externalKey]
		log := ilog.ContainedLog{
			Hash: baseFileHashs[0].Hash,
		}
		for _, fileHash := range externalM[externalKey] {
			files++
			bytes += fileHash.Size
			log.AddRelativePath(fileHash.RelativePath)
			relPaths = append(relPaths, fileHash.RelativePath)
		}
		logBuffer.Append(log)
		externalM.Remove(externalKey)
	}
	return files, bytes, relPaths
}

func removeFiles(basePath string, relativePaths []string) {
	for _, relativePath := range relativePaths {
		path := filepath.Join(basePath, relativePath)
		err := os.Remove(path)
		if err != nil {
			log.Fatal("could not remove contained or duplicate file: ", err)
		}
	}
}
