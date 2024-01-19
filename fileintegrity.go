package fileintegrity

import (
	"github.com/aicirt2012/fileintegrity/src/store"
	"github.com/aicirt2012/fileintegrity/src/store/check"
	"github.com/aicirt2012/fileintegrity/src/store/ilog"
)

// Set with linker flags
var Version = "development"

// Upsert inserts or updates entries into the integrity file. An update is performed when the actual file
// modification date is after the file modification date of the stored entry.
func Upsert(path string, options Options) error {
	return store.Upsert(path, options.toStoreOptions())
}

// Verify verifies that the actual file hash is similar to the hash stored in the integrity file entry.
func Verify(path string, options Options) error {
	return store.Verify(path, options.toStoreOptions())
}

// CheckDuplicates checks for duplicate files within the integrity file.
func CheckDuplicates(path string, options Options) {
	check.Duplicates(path, options.toStoreOptions())
}

// CheckContained checks if files of an external directory are contained within the integrity file.
// With the optional flag fix, contained and duplicated files are deleted form the external directory.
func CheckContained(path string, externalPath string, fix bool, options Options) {
	check.Contained(path, externalPath, fix, options.toStoreOptions())
}

// CheckStyleIssues checks style issues related to the file system based on the integrity file.
// Check categories are: Directory hierarchy issues, path and directory length issues, naming issues.
func CheckStyleIssues(path string, options Options) {
	check.StyleIssues(path, options.toStoreOptions())
}

// CheckExtensionStats checks the distribution of file extensions based on the file size within the integrity file
func CheckExtensionStats(path string, options Options) {
	check.ExtensionStats(path, options.toStoreOptions())
}

// DefaultOptions for execution
func DefaultOptions() Options {
	return Options{
		LogConsole:  true,
		LogFile:     true,
		Backup:      false,
		ProgressBar: true,
	}
}

// LogOptions for execution
func LogOptions(quiet *bool) Options {
	return Options{
		LogConsole:  !*quiet,
		LogFile:     true,
		Backup:      false,
		ProgressBar: !*quiet,
	}
}

// EnabledOptions for execution
func EnabledOptions() Options {
	return Options{
		LogConsole:  true,
		LogFile:     true,
		Backup:      true,
		ProgressBar: false,
	}
}

// DisabledOptions for execution
func DisabledOptions() Options {
	return Options{
		LogConsole:  false,
		LogFile:     false,
		Backup:      false,
		ProgressBar: false,
	}
}

// Options customizable for execution
type Options struct {
	LogConsole  bool
	LogFile     bool
	Backup      bool
	ProgressBar bool
}

func (o Options) toStoreOptions() store.Options {
	return store.Options{
		Log: ilog.Options{
			Console: o.LogConsole,
			File:    o.LogFile,
		},
		Backup:      o.Backup,
		ProgressBar: o.ProgressBar,
	}
}
