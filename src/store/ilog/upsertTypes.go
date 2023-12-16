package ilog

import (
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

type UpsertOperation string

const (
	NEW    UpsertOperation = "NEW"
	UPDATE UpsertOperation = "UPDATE"
	DELETE UpsertOperation = "DELETE"
	SKIP   UpsertOperation = "SKIP"
)

type UpsertLog struct {
	Created      time.Time
	Operation    UpsertOperation
	RelativePath string
}

func (l UpsertLog) serialize() string {
	a := []string{
		l.Created.Format(TimeFormat),
		string(l.Operation),
		l.RelativePath,
	}
	return strings.Join(a, "  ")
}

func (l UpsertLog) visibleOnConsole() bool {
	return true
}

type UpsertSummary struct {
	ExecutionTime time.Duration
	TotalBytes    int64
	HashedBytes   int64
	SkippedFiles  int64
	NewFiles      int64
	UpdatedFiles  int64
	DeletedFiles  int64
}

func (us *UpsertSummary) AddHashedBytes(bytes int64) {
	us.HashedBytes += bytes
}

func (us UpsertSummary) hashRateInS() uint64 {
	s := us.ExecutionTime.Abs().Seconds()
	if s == 0 {
		return 0
	}
	return uint64(float64(us.HashedBytes) / s)
}

func (us UpsertSummary) serialize() string {
	s := title(Upsert)
	s += line("Execution time:", "%.2f s", us.ExecutionTime.Abs().Seconds())
	s += line("Total size:", "%v", humanize.Bytes(uint64(us.TotalBytes)))
	s += line("Hash rate:", "%v/s", humanize.Bytes(us.hashRateInS()))
	s += line("Skipped files:", "%v", us.SkippedFiles)
	s += line("New files:", "%v", us.NewFiles)
	s += line("Updated files:", "%v", us.UpdatedFiles)
	s += line("Deleted files:", "%v", us.DeletedFiles)
	return s
}

func (us UpsertSummary) visibleOnConsole() bool {
	return true
}
