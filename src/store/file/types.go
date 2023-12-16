package file

import (
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

const EmptyHash = "0000000000000000000000000000000000000000000000000000000000000000" // Empty hash means marked for deletion

type FileHash struct {
	Hash         string    `csv:"hash"`
	Created      time.Time `csv:"created"`
	ModTime      time.Time `csv:"mod"`
	Size         int64     `csv:"size"`
	RelativePath string    `csv:"relativePath"`
}

func (fh *FileHash) Equal(o FileHash) bool {
	return fh.Hash == o.Hash &&
		fh.ModTime.Equal(o.ModTime) &&
		fh.Size == o.Size &&
		fh.RelativePath == o.RelativePath
}

type FileHashs []FileHash

func (fs FileHashs) TotalBytes() (s int64) {
	for _, fileHash := range fs {
		s += fileHash.Size
	}
	return s
}

// Implicitly removes duplicated and deleted entries
func (fs FileHashs) DefragmentedMap() FileHashMap {
	m := make(map[string]FileHash)
	for _, fileHash := range fs {
		existingFileHash, exists := m[fileHash.RelativePath]
		if !exists || exists && fileHash.Created.After(existingFileHash.Created) {
			m[fileHash.RelativePath] = fileHash
		}
	}
	for _, fileHash := range fs {
		if fileHash.Hash == EmptyHash {
			delete(m, fileHash.RelativePath)
		}
	}
	return m
}

func (fs FileHashs) Len() int {
	return len(fs)
}
func (fs FileHashs) Swap(i, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}
func (fs FileHashs) Less(i, j int) bool {
	return strings.Compare(fs[i].RelativePath, fs[j].RelativePath) < 0
}

type FileHashMap map[string]FileHash

func (fhm FileHashMap) Has(relativePath string) bool {
	if _, ok := fhm[relativePath]; ok {
		return true
	}
	return false
}

func (fhm *FileHashMap) Remove(relativePath string) {
	delete(*fhm, relativePath)
}

type fileHashsBuffer struct {
	basePath       string
	fileHashs      FileHashs
	maxBytes       int64
	afterFlushHook func()
}

func (fhb *fileHashsBuffer) Append(fileHash FileHash) {
	fhb.fileHashs = append(fhb.fileHashs, fileHash)
	if fhb.fileHashs.TotalBytes() > fhb.maxBytes {
		fhb.Flush()
	}
}

func (fhb *fileHashsBuffer) Flush() {
	if len(fhb.fileHashs) == 0 {
		return
	}
	Append(fhb.basePath, fhb.fileHashs)
	fhb.fileHashs = FileHashs{}
	fhb.afterFlushHook()
}

func NewFileHashsBuffer(basePath string, maxGBytes int64, afterFlashHook func()) fileHashsBuffer {
	return fileHashsBuffer{
		basePath:       basePath,
		fileHashs:      FileHashs{},
		maxBytes:       maxGBytes * humanize.GByte,
		afterFlushHook: afterFlashHook,
	}
}
