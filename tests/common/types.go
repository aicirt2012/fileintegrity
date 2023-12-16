package common

import (
	"time"

	"github.com/aicirt2012/fileintegrity/src/store/ilog"
)

const integrity string = ".integrity"

type File struct {
	relativePath string
	modTime      time.Time
	content      string
}

func NewFile(relativePath string, modTime string, content string) File {
	return File{
		relativePath: NormalizePath(relativePath),
		modTime:      parseTime(modTime),
		content:      content,
	}
}

type Files []File

func (f *Files) Map() map[string]File {
	m := make(map[string]File)
	for _, i := range *f {
		m[i.relativePath] = i
	}
	return m
}

type FileHash struct {
	hash         string
	created      string
	modTime      string
	size         string
	relativePath string
}

func NewFileHash(hash string, created string, modTime string, size string, relativePath string) FileHash {
	return FileHash{
		hash:         hash,
		created:      created,
		modTime:      modTime,
		size:         size,
		relativePath: NormalizePath(relativePath),
	}
}

type LogBlock struct {
	category      ilog.Category
	hash          string
	relativePaths []string
}

func (l LogBlock) firstLine() string {
	return l.category.ToUpper() + " " + l.hash
}

func NewDuplicateLogBlock(hash string, relativePaths []string) LogBlock {
	return LogBlock{
		category:      ilog.Duplicate,
		hash:          hash,
		relativePaths: relativePaths,
	}
}

func NewContainedLogBlock(hash string, relativePaths []string) LogBlock {
	return LogBlock{
		category:      ilog.Contained,
		hash:          hash,
		relativePaths: relativePaths,
	}
}
