package ilog

import (
	"bytes"
	"log"
	"os"
	"strings"
	"time"
)

const TimeFormat = "060102.150405"
const ext = "log"

type serializable interface {
	serialize() string
	visibleOnConsole() bool
}

type Category string

const (
	Upsert    Category = "upsert"
	Verify    Category = "verify"
	Duplicate Category = "duplicate"
	Contained Category = "contained"
	Style     Category = "style"
)

func (c Category) ToUpper() string {
	return strings.ToUpper(string(c))
}

type flushType string

const (
	manual    flushType = "manual"
	automatic flushType = "automatic"
)

type Options struct {
	Console bool
	File    bool
}

type LogFileBuffer struct {
	filename  string
	logs      []serializable
	flushType flushType
	maxItems  uint64
	options   Options
}

func (lf *LogFileBuffer) Flush() {
	if !lf.options.File || len(lf.logs) == 0 {
		return
	}
	f, err := os.OpenFile(lf.filename, os.O_RDWR|os.O_APPEND, 0644)
	if os.IsNotExist(err) {
		f, err = os.Create(lf.filename)
	}
	defer f.Close()
	if err != nil {
		log.Fatal("could not create or open log file")
	}
	var buffer bytes.Buffer
	for _, log := range lf.logs {
		buffer.WriteString(log.serialize() + "\n")
	}
	_, err = f.Write(buffer.Bytes())
	if err != nil {
		log.Fatal("could not write log file")
	}
	lf.logs = []serializable{}
}

func (lf *LogFileBuffer) RequiresFlush() bool {
	return len(lf.logs) > int(lf.maxItems)
}

func (lf *LogFileBuffer) AutomaticFlush() bool {
	return lf.flushType == automatic
}

func (lf *LogFileBuffer) AppendAll(logs []serializable) *LogFileBuffer {
	for _, log := range logs {
		lf.Append(log)
	}
	return lf
}

func (lf *LogFileBuffer) Append(log serializable) *LogFileBuffer {
	if lf.options.Console && log.visibleOnConsole() {
		println(log.serialize())
	}
	lf.logs = append(lf.logs, log)
	if lf.options.File && lf.AutomaticFlush() && lf.RequiresFlush() {
		lf.Flush()
	}
	return lf
}

func (lf *LogFileBuffer) AppendUpsertLog(operation UpsertOperation, relativePath string) *LogFileBuffer {
	lf.Append(UpsertLog{
		Created:      time.Now(),
		Operation:    operation,
		RelativePath: relativePath,
	})
	return lf
}

func (lf *LogFileBuffer) AppendVerifyLog(status VerifyStatus, relativePath string, reason error) *LogFileBuffer {
	lf.Append(VerifyLog{
		Created:      time.Now(),
		Status:       status,
		RelativePath: relativePath,
		Reason:       reason,
	})
	return lf
}
