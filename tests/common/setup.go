package common

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func initFS(wd string, files []File) ([]File, error) {
	for _, file := range files {
		absoluteFilename := filepath.Join(wd, file.relativePath)
		absoluteDir := filepath.Dir(absoluteFilename)
		if err := os.MkdirAll(absoluteDir, 0644); err != nil {
			return nil, err
		}
		if err := os.WriteFile(absoluteFilename, []byte(file.content), 0644); err != nil {
			return nil, err
		}
		if err := os.Chtimes(absoluteFilename, file.modTime, file.modTime); err != nil {
			return nil, err
		}
	}
	return files, nil
}

func upsertScenarioDir(name string) string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("could not get working dir", err)
	}
	dir := filepath.Join(wd, "cases", name)
	err = os.RemoveAll(dir)
	if err != nil {
		log.Fatal("could not wipe", err)
	}
	err = os.MkdirAll(dir, 0644)
	if err != nil {
		log.Fatal("could not create scenario dir", err)
	}
	return dir
}

func CreateScenario(name string, files Files) (string, Files) {
	dir := upsertScenarioDir(name)
	files, err := initFS(dir, files)
	if err != nil {
		log.Fatal("could not init fs", err)
	}
	return dir, files
}

func listFiles(basePath string) (files Files, err error) {
	e := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if basePath == path {
			return nil
		}
		if info.IsDir() && info.Name() == integrity {
			return filepath.SkipDir
		}
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		relativeFilename, err := filepath.Rel(basePath, path)
		if err != nil {
			return err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		files = append(files, File{
			relativePath: NormalizePath(relativeFilename),
			modTime:      info.ModTime(),
			content:      string(content),
		})
		return nil
	})
	return files, e
}

func RemoveFile(dir string, filename string) {
	err := os.Remove(filepath.Join(dir, NormalizePath(filename)))
	if err != nil {
		log.Fatal("could not remove file", err)
	}
}

func UpdateFile(dir string, filename string, content string, modTime string) {
	filename = filepath.Join(dir, NormalizePath(filename))
	err := os.WriteFile(filename, []byte(content), os.ModeAppend)
	if err != nil {
		log.Fatal("could not update file", err)
	}
	mt := parseTime(modTime)
	err = os.Chtimes(filename, mt, mt)
	if err != nil {
		log.Fatal("could not adapt time", err)
	}
}

func lastLogFileContent(dir string) (content string, err error) {
	absoluteDir := filepath.Join(dir, integrity)
	files, err := os.ReadDir(absoluteDir)
	if err != nil {
		return "", err
	}
	for i, j := 0, len(files)-1; i < j; i, j = i+1, j-1 {
		files[i], files[j] = files[j], files[i]
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".log" {
			content, err := os.ReadFile(filepath.Join(absoluteDir, f.Name()))
			return string(content), err
		}
	}
	return "", nil
}

func CreateIntegrityFile(t *testing.T, dir string, fileHashes []FileHash) {
	absoluteDir := filepath.Join(dir, integrity)
	err := os.Mkdir(absoluteDir, os.ModeAppend)
	if err != nil {
		log.Fatal("could not create integrity dir", err)
	}
	absoluteFilename := filepath.Join(absoluteDir, integrity)
	content := ""
	for _, fh := range fileHashes {
		arr := []string{fh.hash, fh.created, fh.modTime, fh.size, fh.relativePath}
		content += strings.Join(arr, ",") + "\n"
	}
	err = os.WriteFile(absoluteFilename, []byte(content), 0644)
	if err != nil {
		log.Fatal("could not create integrity file", err)
	}
}

func parseTime(value string) time.Time {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		log.Fatal("could not parse test data modTime ", value)
	}
	return t
}
