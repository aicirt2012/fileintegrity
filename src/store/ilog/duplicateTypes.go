package ilog

import (
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

type DuplicateLog struct {
	Hash          string
	RelativePaths []string
}

func (l *DuplicateLog) AddRelativePath(path string) {
	l.RelativePaths = append(l.RelativePaths, path)
}

func (l DuplicateLog) serialize() string {
	var b strings.Builder
	b.WriteString(Duplicate.ToUpper() + " " + l.Hash + "\n")
	b.WriteString(strings.Join(l.RelativePaths, "\n"))
	b.WriteString("\n")
	return b.String()
}

func (l DuplicateLog) visibleOnConsole() bool {
	return true
}

type DuplicateSummary struct {
	ExecutionTime  time.Duration
	TotalFiles     int64
	TotalBytes     int64
	DuplicateFiles int64
	DuplicateBytes int64
}

func (ds DuplicateSummary) filePercentage() float64 {
	if ds.DuplicateFiles == 0 {
		return 0
	}
	return float64(ds.DuplicateFiles) / float64(ds.TotalFiles) * 100
}

func (ds DuplicateSummary) bytePercentage() float64 {
	if ds.DuplicateBytes == 0 {
		return 0
	}
	return float64(ds.DuplicateBytes) / float64(ds.TotalBytes) * 100
}

func (ds DuplicateSummary) serialize() string {
	s := title(Duplicate)
	s += line("Execution time:", "%.2f s", ds.ExecutionTime.Abs().Seconds())
	s += line("Total files:", "%v", ds.TotalFiles)
	s += line("Duplicate files:", "%v", ds.DuplicateFiles)
	s += line("Duplicate file percentage:", "%.1f", ds.filePercentage())
	s += line("Total size:", "%v", humanize.Bytes(uint64(ds.TotalBytes)))
	s += line("Duplicate size:", "%v", humanize.Bytes(uint64(ds.DuplicateBytes)))
	s += line("Duplicate size percentage:", "%.1f", ds.bytePercentage())
	return s
}

func (ds DuplicateSummary) visibleOnConsole() bool {
	return true
}
