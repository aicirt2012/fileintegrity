package common

import (
	"sort"

	"github.com/aicirt2012/fileintegrity/src/store/ilog"
	"golang.org/x/exp/maps"
)

type LogStyleMap map[string]ilog.StyleLog

func (m LogStyleMap) Put(log ilog.StyleLog) {
	m[log.String()] = log
}

func (m LogStyleMap) PutAll(logs []ilog.StyleLog) {
	for _, log := range logs {
		m.Put(log)
	}
}

func (m LogStyleMap) SortValues() []ilog.StyleLog {
	logs := maps.Values(m)
	sort.Slice(logs, func(i, j int) bool {
		return logs[i].RelativePath < logs[j].RelativePath
	})
	return logs
}

func (m LogStyleMap) WriteTo(logBuffer *ilog.LogFileBuffer) int {
	for _, log := range m.SortValues() {
		logBuffer.Append(log)
	}
	return len(m)
}
