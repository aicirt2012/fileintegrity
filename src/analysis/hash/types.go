package hash

import "time"

type VerifyRequest struct {
	BasePath     string
	RelativePath string
	Size         int64
	ModTime      time.Time
	Hash         string
}

type VerifyResponse struct {
	RelativePath string
	Error        error
}

type CreateRequest struct {
	BasePath     string
	RelativePath string
}

type CreateResponse struct {
	RelativePath string
	Hash         string
	Error        error
}
