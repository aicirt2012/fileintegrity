package dir

import (
	"log"
	"os"
	"path/filepath"
)

const Name = ".integrity"

func AssertDir(path string) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Fatal("asserted dir does not exist", path)
	} else if err != nil {
		log.Fatal("assert dir fails", err)
	}
	if !info.IsDir() {
		log.Fatal("asserted dir is not a directory", path)
	}
}

func AssertIntegrityDir(basePath string) {
	AssertDir(filepath.Join(basePath, Name))
}

func UpsertIntegrityDir(basePath string) {
	path := filepath.Join(basePath, Name)
	if info, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, 0644); err != nil {
			log.Fatal("could not create integrity dir", err)
		}
		if err = hideFile(path); err != nil {
			log.Fatal("could not hide integrity dir", err)
		}
	} else if err != nil || !info.IsDir() {
		log.Fatal("could not ensure integrity dir", err)
	}
}
