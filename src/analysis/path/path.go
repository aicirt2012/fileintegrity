package path

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"

	"github.com/aicirt2012/fileintegrity/src/store/dir"

	"golang.org/x/exp/slices"
)

var ignoredDirs = []string{
	"$RECYCLE.BIN",              // win
	"System Volume Information", // win
	".Spotlight-V100",           // macos
	".fseventsd",                // macos
	".Trashes",                  // macos
	".TemporaryItems",           // macos
	dir.Name,                    // integrity store

}

var ignoredFiles = []*regexp.Regexp{
	regexp.MustCompile(`^~\$.*$`),      // win ~$*
	regexp.MustCompile(`^\._.*$`),      // macos ._*
	regexp.MustCompile(`^\.DS_Store$`), // macos
}

func ComputeDiskFileMap(basePath string) (DiskFileMap, error) {
	diskFileMap := DiskFileMap{}
	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if basePath == path {
			return nil
		}
		if info.IsDir() && isIgnoreDir(info.Name()) {
			return filepath.SkipDir
		}
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if isIgnoredFile(info.Name()) {
			return nil
		}
		relPath, err := filepath.Rel(basePath, path)
		if err != nil {
			return errors.New("could not extract relative path from: " + path)
		}
		diskFileMap.Add(DiskFile{
			AbsolutePath: path,
			RelativePath: relPath,
			Size:         info.Size(),
			ModTime:      info.ModTime(),
		})
		return nil
	})
	return diskFileMap, err
}

func isIgnoreDir(name string) bool {
	return slices.Contains(ignoredDirs, name)
}

func isIgnoredFile(name string) bool {
	for _, pattern := range ignoredFiles {
		if pattern.MatchString(name) {
			return true
		}
	}
	return false
}
