package license

import (
	"embed"
	"errors"
)

//go:embed license.txt
var fs embed.FS

func Text() (string, error) {
	content, err := fs.ReadFile("license.txt")
	if err != nil {
		return "", errors.New("no license file embedded")
	}
	return string(content), nil
}
