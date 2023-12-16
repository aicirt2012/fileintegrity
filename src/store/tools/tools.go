package tools

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aicirt2012/fileintegrity/src/store"
	"github.com/aicirt2012/fileintegrity/src/store/file"
	"github.com/aicirt2012/fileintegrity/src/store/ilog"
	"golang.org/x/exp/maps"
)

func FindDuplicates(basePath string, options store.Options) {
	start := time.Now()
	summary := ilog.DuplicateSummary{}
	logBuffer := ilog.NewAutomaticLogBuffer(basePath, ilog.Duplicate, 1000, options.Log)

	m, tf, tb := calcHashSizeMap(file.LoadContent(basePath))
	summary.TotalFiles = tf
	summary.TotalBytes = tb

	df, db, _ := analyzeDuplicates(m, &logBuffer)
	summary.DuplicateFiles = df
	summary.DuplicateBytes = db

	summary.ExecutionTime = time.Since(start)
	logBuffer.Append(summary).Flush()
}

func FindContained(basePath string, externalPath string, delete bool, options store.Options) {
	start := time.Now()
	summary := ilog.ContainedSummary{}
	logBuffer := ilog.NewAutomaticLogBuffer(basePath, ilog.Contained, 1000, options.Log)

	options.Backup = true
	err := store.Upsert(externalPath, options)
	if err != nil {
		log.Fatal(err)
	}
	baseM, _, _ := calcHashSizeMap(file.LoadContent(basePath))
	externalM, tf, tb := calcHashSizeMap(file.LoadContent(externalPath))
	summary.TotalFiles = tf
	summary.TotalBytes = tb

	cf, cb, containedRelPaths := analyzeContained(baseM, externalM, &logBuffer)
	summary.ContainedFiles = cf
	summary.ContainedBytes = cb

	df, db, duplicateRelPaths := analyzeDuplicates(externalM, &logBuffer)
	summary.DuplicateFiles = df
	summary.DuplicateBytes = db

	if delete {
		removeFiles(externalPath, append(containedRelPaths, duplicateRelPaths...))
	}

	summary.ExecutionTime = time.Since(start)
	logBuffer.Append(summary).Flush()
}

func analyzeDuplicates(m uniqueMap, logBuffer *ilog.LogFileBuffer) (int64, int64, []string) {
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

func analyzeContained(baseM uniqueMap, externalM uniqueMap, logBuffer *ilog.LogFileBuffer) (int64, int64, []string) {
	var relPaths []string
	var files, bytes int64
	for _, externalKey := range maps.Keys(externalM) {
		if !baseM.has(externalKey) {
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
		externalM.remove(externalKey)
	}
	return files, bytes, relPaths
}

// Create map with hash and size as key, ignore small files as well as git files
func calcHashSizeMap(fileHashs []file.FileHash) (uniqueMap, int64, int64) {
	m := uniqueMap{}
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

func removeFiles(basePath string, relativePaths []string) {
	for _, relativePath := range relativePaths {
		path := filepath.Join(basePath, relativePath)
		err := os.Remove(path)
		if err != nil {
			log.Fatal("could not remove contained file: ", err)
		}
	}
}
