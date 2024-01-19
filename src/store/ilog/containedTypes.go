package ilog

import (
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

type ContainedLog struct {
	Hash          string
	RelativePaths []string
}

func (l *ContainedLog) AddRelativePath(path string) {
	l.RelativePaths = append(l.RelativePaths, path)
}

func (l ContainedLog) serialize() string {
	var b strings.Builder
	b.WriteString(Contains.ToUpper() + " " + l.Hash + "\n")
	b.WriteString(strings.Join(l.RelativePaths, "\n"))
	b.WriteString("\n")
	return b.String()
}

func (l ContainedLog) visibleOnConsole() bool {
	return true
}

type ContainedSummary struct {
	ExecutionTime  time.Duration
	TotalFiles     int64
	TotalBytes     int64
	ContainedFiles int64
	ContainedBytes int64
	DuplicateFiles int64
	DuplicateBytes int64
}

func (ds ContainedSummary) overheadFiles() int64 {
	return ds.ContainedFiles + ds.DuplicateFiles
}

func (ds ContainedSummary) overheadBytes() int64 {
	return ds.ContainedBytes + ds.DuplicateBytes
}

func (ds ContainedSummary) overheadFilePercentage() float64 {
	if ds.overheadFiles() == 0 {
		return 0
	}
	return float64(ds.overheadFiles()) / float64(ds.TotalFiles) * 100
}

func (ds ContainedSummary) overheadBytePercentage() float64 {
	if ds.overheadBytes() == 0 {
		return 0
	}
	return float64(ds.overheadBytes()) / float64(ds.TotalBytes) * 100
}

func (ds ContainedSummary) serialize() string {
	s := title(Contains)
	s += line("Execution time:", "%.2f s", ds.ExecutionTime.Abs().Seconds())
	s += line("Total files:", "%v", ds.TotalFiles)
	s += line("Contained files:", "%v", ds.ContainedFiles)
	s += line("Duplicate files:", "%v", ds.DuplicateFiles)
	s += line("Overhead file percentage:", "%.1f", ds.overheadFilePercentage())
	s += line("Total size:", "%v", humanize.Bytes(uint64(ds.TotalBytes)))
	s += line("Contained size:", "%v", humanize.Bytes(uint64(ds.ContainedBytes)))
	s += line("Duplicate size:", "%v", humanize.Bytes(uint64(ds.DuplicateBytes)))
	s += line("Overhead size percentage:", "%.1f", ds.overheadBytePercentage())
	return s
}

func (ds ContainedSummary) visibleOnConsole() bool {
	return true
}
