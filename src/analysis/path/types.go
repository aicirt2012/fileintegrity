package path

import (
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type DiskFileMap map[string]DiskFile

func (dfm DiskFileMap) Has(path string) bool {
	if _, ok := dfm[path]; ok {
		return true
	}
	return false
}

func (dfm DiskFileMap) Size() int {
	return len(dfm)
}

func (dfm DiskFileMap) TotalBytes() int64 {
	sum := int64(0)
	for _, diskFile := range dfm {
		sum += diskFile.Size
	}
	return sum
}

func (dfm DiskFileMap) Add(p DiskFile) {
	dfm[p.RelativePath] = p
}

func (dfm *DiskFileMap) Remove(relativePath string) {
	delete(*dfm, relativePath)
}

func (dfm DiskFileMap) String() string {
	s := "\nDiskFileMap\n"
	for _, diskFile := range dfm {
		s += diskFile.String() + "\n"
	}
	s += "---\n"
	return s
}

func (dfm DiskFileMap) Print() {
	println(dfm.String())
}

type DiskFile struct {
	AbsolutePath string
	RelativePath string
	Size         int64
	ModTime      time.Time
}

func (df DiskFile) Base() string {
	return filepath.Clean(strings.TrimSuffix(df.AbsolutePath, df.RelativePath))
}

func (df DiskFile) Equals(o DiskFile) bool {
	return df.ModTime == o.ModTime && df.Size == o.Size
}

func (df DiskFile) unixNano() string {
	return strconv.FormatInt(df.ModTime.UTC().UnixNano(), 10)
}

func (df DiskFile) String() string {
	s := df.unixNano() + "  "
	s += strconv.FormatInt(df.Size, 10) + "  "
	s += df.RelativePath + "\n"
	return s
}
