package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
)

func CreationWorker(requests <-chan CreateRequest, responses chan<- CreateResponse) {
	for request := range requests {
		hash, err := Hash(filepath.Join(request.BasePath, request.RelativePath))
		responses <- CreateResponse{
			RelativePath: request.RelativePath,
			Hash:         hash,
			Error:        err,
		}
	}
}

func VerifyWorker(requests <-chan VerifyRequest, responses chan<- VerifyResponse) {
	for request := range requests {
		responses <- VerifyResponse{
			RelativePath: request.RelativePath,
			Error:        verify(request),
		}
	}
}

func Hash(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", errors.New("Could not open file for hashing: " + filename)
	}
	defer file.Close()
	buf := make([]byte, 30*1024*1024)
	sha256 := sha256.New()
	for {
		n, err := file.Read(buf)
		if n > 0 {
			_, err := sha256.Write(buf[:n])
			if err != nil {
				log.Fatal(err)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Read %d bytes: %v", n, err)
			return "", errors.New("error during hashing")
		}
	}
	return hex.EncodeToString(sha256.Sum(nil)), nil
}

func verify(request VerifyRequest) error {
	path := filepath.Join(request.BasePath, request.RelativePath)
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return errors.New("file does not exist")
	} else if err != nil {
		return errors.New("file stat unreadable")
	}
	if info.Size() != request.Size {
		return errors.New("file size different")
	}
	hash, err := Hash(path)
	if err != nil {
		return err
	}
	if hash != request.Hash {
		return errors.New("file hash different")
	}
	return nil
}
