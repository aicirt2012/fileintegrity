package store

import (
	"runtime"
	"time"

	"github.com/aicirt2012/fileintegrity/src/analysis/hash"
	"github.com/aicirt2012/fileintegrity/src/analysis/path"
	"github.com/aicirt2012/fileintegrity/src/store/dir"
	"github.com/aicirt2012/fileintegrity/src/store/file"
	"github.com/aicirt2012/fileintegrity/src/store/ilog"

	"golang.org/x/exp/maps"
)

func Upsert(basePath string, options Options) error {
	dir.AssertDir(basePath)
	dir.UpsertIntegrityDir(basePath)
	if options.Backup {
		file.Backup(basePath)
	}
	start := time.Now()

	logBuffer := ilog.NewManualLogBuffer(basePath, ilog.Upsert, options.Log)
	fileBuffer := file.NewFileHashsBuffer(basePath, 1, logBuffer.Flush)

	fileHashes := file.LoadContent(basePath)
	fileHashMap := fileHashes.DefragmentedMap()
	diskFileMap, err := path.ComputeDiskFileMap(basePath)
	if err != nil {
		return err
	}

	summary := ilog.UpsertSummary{
		TotalBytes: diskFileMap.TotalBytes(),
	}

	progressBar := ilog.ProgressBar(summary.TotalBytes, options.ProgressBar)

	// Process unchanged entries
	for _, file := range maps.Values(diskFileMap) {
		fileHash, exists := fileHashMap[file.RelativePath]
		if exists && fileHash.ModTime.Equal(file.ModTime) && fileHash.Size == file.Size {
			diskFileMap.Remove(fileHash.RelativePath)
			fileHashMap.Remove(fileHash.RelativePath)
			progressBar.Add64(file.Size)
			summary.SkippedFiles++
		}
	}

	// Initialize channels
	requests := make(chan hash.CreateRequest, 10)
	responses := make(chan hash.CreateResponse, 100)
	await := make(chan bool)

	// Consume file hash responses
	go func() {
		for i := 0; i < len(diskFileMap); i++ {
			response := <-responses
			diskFile := diskFileMap[response.RelativePath]
			fileBuffer.Append(file.FileHash{
				Hash:         response.Hash,
				Created:      time.Now(),
				ModTime:      diskFile.ModTime,
				Size:         diskFile.Size,
				RelativePath: response.RelativePath,
			})
			if fileHashMap.Has(response.RelativePath) {
				logBuffer.AppendUpsertLog(ilog.UPDATE, response.RelativePath)
				summary.UpdatedFiles++
			} else {
				logBuffer.AppendUpsertLog(ilog.NEW, response.RelativePath)
				summary.NewFiles++
			}
			summary.AddHashedBytes(diskFile.Size)
			progressBar.Add64(diskFile.Size)
		}
		fileBuffer.Flush()
		await <- true
	}()

	// Create file hash workers
	for w := 1; w <= runtime.NumCPU(); w++ {
		go hash.CreationWorker(requests, responses)
	}

	// Produce file hash requests for all new or not up to date entries
	for _, file := range maps.Values(diskFileMap) {
		requests <- hash.CreateRequest{
			BasePath:     basePath,
			RelativePath: file.RelativePath,
		}
	}
	<-await

	// Delete hashes for non existing files
	for _, hash := range maps.Values(fileHashMap) {
		if _, exists := diskFileMap[hash.RelativePath]; exists {
			continue
		}
		fileBuffer.Append(file.FileHash{
			Hash:         file.EmptyHash,
			Created:      time.Now(),
			ModTime:      hash.ModTime,
			Size:         hash.Size,
			RelativePath: hash.RelativePath,
		})
		logBuffer.AppendUpsertLog(ilog.DELETE, hash.RelativePath)
		summary.DeletedFiles++
	}

	fileBuffer.Flush()
	file.Defragment(basePath)
	summary.ExecutionTime = time.Since(start)
	logBuffer.Append(summary).Flush()

	return nil
}

func Verify(basePath string, options Options) error {
	dir.AssertDir(basePath)
	dir.AssertIntegrityDir(basePath)
	start := time.Now()
	logBuffer := ilog.NewAutomaticLogBuffer(basePath, ilog.Verify, 1000, options.Log)
	fileHashes := file.LoadContent(basePath)
	fileHashesMap := fileHashes.DefragmentedMap()
	totalBytes := fileHashes.TotalBytes()
	errorCount := int64(0)

	progressBar := ilog.ProgressBar(totalBytes, options.ProgressBar)

	// Initialize channels
	requests := make(chan hash.VerifyRequest, 10)
	responses := make(chan hash.VerifyResponse, 100)
	await := make(chan bool)

	// Consume file hash responses
	go func() {
		for i := 0; i < len(fileHashes); i++ {
			response := <-responses
			fileHash := fileHashesMap[response.RelativePath]

			status := ilog.OK
			var reason error
			if response.Error != nil {
				status = ilog.ERROR
				reason = response.Error
				errorCount++
			}
			logBuffer.AppendVerifyLog(status, response.RelativePath, reason)
			progressBar.Add64(fileHash.Size)
		}
		await <- true
	}()

	// Create file hash workers
	for w := 1; w <= runtime.NumCPU(); w++ {
		go hash.VerifyWorker(requests, responses)
	}

	// Produce file hash requests
	for _, fileHash := range fileHashes {
		requests <- hash.VerifyRequest{
			BasePath:     basePath,
			RelativePath: fileHash.RelativePath,
			Size:         fileHash.Size,
			ModTime:      fileHash.ModTime,
			Hash:         fileHash.Hash,
		}
	}
	<-await

	logBuffer.Append(ilog.VerifySummary{
		ExecutionTime: time.Since(start),
		TotalBytes:    totalBytes,
		ValidFiles:    int64(len(fileHashes) - int(errorCount)),
		InvalidFiles:  errorCount,
	}).Flush()
	return nil
}
