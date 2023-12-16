package ilog

import (
	"strings"
	"time"
)

type IssueType string

const (
	HIERARCHY_ISSUE IssueType = "HIERARCHY ISSUE"
	NAMING_ISSUE    IssueType = "NAMING ISSUE"
	LENGTH_ISSUE    IssueType = "LENGTH ISSUE"
)

type StyleLog struct {
	IssueType    IssueType
	Reason       string
	RelativePath string
}

func (l StyleLog) String() string {
	return l.RelativePath + l.Reason
}

func (l StyleLog) serialize() string {
	lines := []string{
		string(l.IssueType) + ": " + l.Reason,
		l.RelativePath,
	}
	return strings.Join(lines, "\n") + "\n"
}

func (l StyleLog) visibleOnConsole() bool {
	return true
}

type StyleSummary struct {
	ExecutionTime   time.Duration
	HierarchyIssues int
	NamingIssues    int
	LengthIssues    int
	TotalDirs       int64
}

func (ds StyleSummary) serialize() string {
	s := title(Style)
	s += line("Execution time:", "%.2f s", ds.ExecutionTime.Abs().Seconds())
	s += line("Hierarchy issues:", "%v", ds.HierarchyIssues)
	s += line("Naming issues:", "%v", ds.NamingIssues)
	s += line("Length issues:", "%v", ds.LengthIssues)
	return s
}

func (ds StyleSummary) visibleOnConsole() bool {
	return true
}
