package ilog

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/aicirt2012/fileintegrity/src/store/dir"
)

const summaryColumns = 42

func NewManualLogBuffer(basePath string, category Category, options Options) LogFileBuffer {
	return LogFileBuffer{
		filename:  generateFilename(basePath, category),
		flushType: manual,
		options:   options,
	}
}

func NewAutomaticLogBuffer(basePath string, category Category, maxItems uint64, options Options) LogFileBuffer {
	return LogFileBuffer{
		filename:  generateFilename(basePath, category),
		maxItems:  maxItems,
		flushType: automatic,
		options:   options,
	}
}

func ProgressBar(max int64, visible bool) *progressbar.ProgressBar {
	return progressbar.NewOptions64(
		max,
		progressbar.OptionClearOnFinish(),
		progressbar.OptionThrottle(300*time.Microsecond),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetVisibility(visible),
	)
}

func generateFilename(basePath string, category Category) string {
	name := strings.Join([]string{
		time.Now().Format(TimeFormat),
		string(category),
		ext,
	}, ".")
	return filepath.Join(basePath, dir.Name, name)
}

func title(c Category) string {
	category := cases.Title(language.AmericanEnglish).String(string(c))
	left := fmt.Sprintf("//// %v Summary ", category)
	line := padRight(left, "/", summaryColumns)
	return "\n\n" + line + "\n"
}

func line(label string, valueFmt string, value interface{}) string {
	fmtValue := fmt.Sprintf(valueFmt, value)
	line := padMiddle(label, fmtValue, " ", summaryColumns)
	return line + "\n"
}

func padRight(s string, padStr string, n int) string {
	i := n - len(s)
	pad := ""
	if i > 0 {
		pad = strings.Repeat(padStr, i)
	}
	return s + pad
}

func padMiddle(left string, right string, padStr string, n int) string {
	i := n - len(left+right)
	pad := ""
	if i > 0 {
		pad = strings.Repeat(padStr, i)
	}
	return left + pad + right
}
