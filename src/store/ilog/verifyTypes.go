package ilog

import (
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

type VerifyStatus string

const (
	OK    VerifyStatus = "OK"
	ERROR VerifyStatus = "ERROR"
)

type VerifyLog struct {
	Created      time.Time
	Status       VerifyStatus
	RelativePath string
	Reason       error
}

func (l VerifyLog) serialize() string {
	reason := ""
	if l.Reason != nil {
		reason = l.Reason.Error()
	}
	a := []string{
		l.Created.Format(TimeFormat),
		string(l.Status),
		l.RelativePath,
		reason,
	}
	return strings.Join(a, "  ")
}

func (l VerifyLog) visibleOnConsole() bool {
	return l.Status == ERROR
}

type VerifySummary struct {
	ExecutionTime time.Duration
	TotalBytes    int64
	ValidFiles    int64
	InvalidFiles  int64
}

func (vs VerifySummary) invalidFilesPercentage() float64 {
	return float64(vs.InvalidFiles) / float64((vs.InvalidFiles + vs.ValidFiles)) * 100
}

func (vs VerifySummary) hashRateInS() uint64 {
	s := vs.ExecutionTime.Abs().Seconds()
	if s == 0 {
		return 0
	}
	return uint64(float64(vs.TotalBytes) / s)
}

func (l VerifySummary) serialize() string {
	s := title(Verify)
	s += line("Execution time:", "%.2f s", l.ExecutionTime.Abs().Seconds())
	s += line("Total size:", "%v", humanize.Bytes(uint64(l.TotalBytes)))
	s += line("Hash rate:", "%v/s", humanize.Bytes(l.hashRateInS()))
	s += line("Verified valid files:", "%v", l.ValidFiles)
	s += line("Verified invalid files:", "%v", l.InvalidFiles)
	s += line("Percentage of invalid files:", "%.6f", l.invalidFilesPercentage())
	return s
}

func (l VerifySummary) visibleOnConsole() bool {
	return true
}
