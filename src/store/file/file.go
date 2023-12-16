package file

import (
	"archive/zip"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/aicirt2012/fileintegrity/src/store/dir"
	"github.com/aicirt2012/fileintegrity/src/store/ilog"
	"github.com/mitchellh/hashstructure/v2"

	"github.com/gocarina/gocsv"
	"golang.org/x/exp/maps"
)

const name string = ".integrity"

var mu sync.Mutex

func LoadContent(basePath string) FileHashs {
	return loadContentInternal(basePath, true)
}

func Append(basePath string, fileHashs FileHashs) {
	mu.Lock()
	defer mu.Unlock()
	f := openOrCreateFile(basePath)
	defer f.Close()
	line, err := gocsv.MarshalStringWithoutHeaders(&fileHashs)
	if err != nil {
		log.Fatal("could not serialize file hash", err)
	}
	_, err = f.WriteString(line)
	if err != nil {
		log.Fatal("could not write integrity file", err)
	}
}

// During execution new hashes are only appended in the integrity file due to performance reasons.
// This may leads to duplicate entries which are eliminated in a final step.
func Defragment(basePath string) {
	mu.Lock()
	defer mu.Unlock()
	filename := filepath.Join(basePath, dir.Name, name)
	fileHashs := loadContentInternal(basePath, false)
	hashsBefore := hash(fileHashs)

	// Remove duplicated and deleted file entries
	m := fileHashs.DefragmentedMap()
	uniqueFileHashes := maps.Values(m)
	sort.Sort(FileHashs(uniqueFileHashes))

	// Detect unchanged content to prevent change of modification date
	if hashsBefore == hash(uniqueFileHashes) {
		return
	}

	// Override existing file with new serialized content
	integrityFile, err := os.Create(filename)
	if err != nil {
		log.Fatal("could not open integrity file", err)
	}
	defer integrityFile.Close()

	err = gocsv.MarshalWithoutHeaders(&uniqueFileHashes, integrityFile)
	if err != nil {
		log.Fatal("could not serialize integrity file", err)
	}
}

func Backup(basePath string) {
	mu.Lock()
	defer mu.Unlock()
	integrityFilename := filepath.Join(basePath, dir.Name, name)
	integrityInfo, err := os.Stat(integrityFilename)
	if os.IsNotExist(err) {
		return // if not exist, an backup is not required
	}
	integrityFileContent, err := os.ReadFile(integrityFilename)
	if err != nil {
		log.Fatal("could not backup integrity file", err)
	}

	zipFilename := filepath.Join(basePath, dir.Name,
		integrityInfo.ModTime().Format(ilog.TimeFormat)+name+".zip")
	zipFile, err := os.Create(zipFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer zipFile.Close()

	wr := zip.NewWriter(zipFile)
	defer wr.Close()

	f, err := wr.Create(name)
	if err != nil {
		log.Fatal("could not add integrity file to zip", err)
	}

	_, err = f.Write(integrityFileContent)
	if err != nil {
		log.Fatal("could not write data to the integrity zip file", err)
	}
}

func loadContentInternal(basePath string, lock bool) FileHashs {
	if lock {
		mu.Lock()
		defer mu.Unlock()
	}
	fileHashes := FileHashs{}
	f := openOrCreateFile(basePath)
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		log.Fatal("could not get stats of integrity file", err)
	}
	if info.Size() > 0 {
		err = gocsv.UnmarshalWithoutHeaders(f, &fileHashes)
		if err != nil {
			log.Fatal("could not deserialize integrity file", err)
		}
	}
	return fileHashes
}

func openOrCreateFile(basePath string) *os.File {
	filename := filepath.Join(basePath, dir.Name, name)
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0644)
	if os.IsNotExist(err) {
		f, err = os.Create(filename)
	}
	if err != nil {
		log.Fatal("could not create or open integrity file", err)
	}
	return f
}

func hash(i interface{}) uint64 {
	hash, err := hashstructure.Hash(i, hashstructure.FormatV2, nil)
	if err != nil {
		log.Fatal("could not hash", err)
	}
	return hash
}
