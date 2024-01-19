package ilog

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
)

type ExtensionStatsLog struct {
	Name       string
	Bytes      int64
	Percentage float64
}

func (l ExtensionStatsLog) serialize() string {
	p := fmt.Sprintf("%.3f%%", l.Percentage)
	b := humanize.Bytes(uint64(l.Bytes))
	return fmt.Sprintf("%7s  %6s  *%v", p, b, l.Name)
}

func (l ExtensionStatsLog) visibleOnConsole() bool {
	return true
}

type ExtensionStatsSummary struct {
	ExecutionTime    time.Duration
	TotalBytes       int64
	UniqueExtensions int
}

func (ds ExtensionStatsSummary) serialize() string {
	s := title(ExtensionStats)
	s += line("Execution time:", "%.2f s", ds.ExecutionTime.Abs().Seconds())
	s += line("Total size:", "%v", humanize.Bytes(uint64(ds.TotalBytes)))
	s += line("Unique extensions:", "%v", ds.UniqueExtensions)
	return s
}

func (ds ExtensionStatsSummary) visibleOnConsole() bool {
	return true
}
