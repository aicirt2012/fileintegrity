package store

import "github.com/aicirt2012/fileintegrity/src/store/ilog"

type Options struct {
	Log         ilog.Options
	Backup      bool
	ProgressBar bool
}
