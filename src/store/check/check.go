package check

import (
	"github.com/aicirt2012/fileintegrity/src/store"
	"github.com/aicirt2012/fileintegrity/src/store/check/contained"
	"github.com/aicirt2012/fileintegrity/src/store/check/duplicate"
	"github.com/aicirt2012/fileintegrity/src/store/check/style"
)

func Contained(basePath string, externalPath string, fix bool, options store.Options) {
	contained.Check(basePath, externalPath, fix, options)
}

func Duplicates(basePath string, options store.Options) {
	duplicate.Check(basePath, options)
}

func StyleIssues(basePath string, options store.Options) {
	style.Check(basePath, options)
}
