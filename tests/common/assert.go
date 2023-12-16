package common

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/aicirt2012/fileintegrity/src/store/ilog"
	"github.com/stretchr/testify/assert"
)

func AssertFilesExist(t *testing.T, basePath string, expectedFiles Files) {
	actualFiles, err := listFiles(basePath)
	if err != nil {
		assert.Fail(t, "could not list actual files", err)
	}
	for _, ef := range expectedFiles {
		af := actualFiles.Map()[ef.relativePath]
		assert.Equal(t, ef.relativePath, af.relativePath)
		assert.Equal(t, ef.content, af.content)
		assert.Equal(t, ef.modTime.UTC(), af.modTime.UTC())
	}
	assert.Len(t, actualFiles, len(expectedFiles))
}

func AssertIntegrityFile(t *testing.T, dir string, expected []FileHash) {
	filename := filepath.Join(dir, integrity, integrity)
	assert.DirExists(t, filepath.Join(dir, integrity))
	assert.FileExists(t, filename)
	bytes, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal("could not read integrity file", err)
	}
	content := string(bytes)
	actual := [][]string{}
	for _, line := range strings.Split(content, "\n") {
		if line == "" {
			break
		}
		cells := strings.Split(line, ",")
		actual = append(actual, cells)
	}

	assert.Len(t, actual, len(expected))
	for i := range expected {
		expectedLine := expected[i]
		actualLine := actual[i]
		assert.Len(t, actualLine, 5)
		assert.Equal(t, actualLine[0], expectedLine.hash, expectedLine.relativePath)
		assert.Equal(t, parseTime(actualLine[2]).UTC(), parseTime(expectedLine.modTime).UTC())
		assert.Equal(t, actualLine[3], expectedLine.size)
		assert.Equal(t, NormalizePath(actualLine[4]), expectedLine.relativePath)
	}
}

func AssertUpsertLogFile(t *testing.T, dir string, skipped int, new int, updated int, deleted int) {
	expectedLogLines := new + updated + deleted
	content, err := lastLogFileContent(dir)
	if err != nil {
		log.Fatal("could not read log file", err)
	}
	lines := strings.Split(content, "\n")

	regex := regexp.MustCompile(`^\d{6}\.\d{6}  (NEW|UPDATE|DELETE)  .{1,260}$`)
	for i := 0; i < expectedLogLines; i++ {
		assert.Regexp(t, regex, lines[i])
	}

	assert.Equal(t, "//// Upsert Summary //////////////////////", lines[expectedLogLines+2])
	assertSummaryLine(t, "Skipped files:", skipped, lines[expectedLogLines+6])
	assertSummaryLine(t, "New files:", new, lines[expectedLogLines+7])
	assertSummaryLine(t, "Updated files:", updated, lines[expectedLogLines+8])
	assertSummaryLine(t, "Deleted files:", deleted, lines[expectedLogLines+9])
}

func AssertVerifyLogFile(t *testing.T, dir string, valid int, invalid int) {
	expectedLogLines := valid + invalid
	content, err := lastLogFileContent(dir)
	if err != nil {
		log.Fatal("could not read log file", err)
	}
	lines := strings.Split(content, "\n")

	regex := regexp.MustCompile(`^\d{6}\.\d{6}  (OK|ERROR)  .{1,260}$`)
	for i := 0; i < expectedLogLines; i++ {
		assert.Regexp(t, regex, lines[i])
	}

	assert.Equal(t, "//// Verify Summary //////////////////////", lines[expectedLogLines+2])
	assertSummaryLine(t, "Verified valid files:", valid, lines[expectedLogLines+6])
	assertSummaryLine(t, "Verified invalid files:", invalid, lines[expectedLogLines+7])
}

func AssertDuplicateLogFile(t *testing.T, dir string, blocks []LogBlock, unique int, duplicates int) {
	content, err := lastLogFileContent(dir)
	if err != nil {
		log.Fatal("could not read log file", err)
	}
	lines := strings.Split(content, "\n")

	currentLine, actualDuplicates, _ := AssertLogBlocks(t, lines, blocks)
	assert.Equal(t, duplicates, actualDuplicates)

	assert.Equal(t, "//// Duplicate Summary ///////////////////", lines[currentLine+2])
	assertSummaryLine(t, "Total files:", unique+duplicates, lines[currentLine+4])
	assertSummaryLine(t, "Duplicate files:", duplicates, lines[currentLine+5])
}

func AssertContainedLogFile(t *testing.T, dir string, blocks []LogBlock, unique int, duplicates int, contained int) {
	content, err := lastLogFileContent(dir)
	if err != nil {
		log.Fatal("could not read log file", err)
	}
	lines := strings.Split(content, "\n")

	currentLine, aDuplicates, aContained := AssertLogBlocks(t, lines, blocks)
	assert.Equal(t, duplicates, aDuplicates)
	assert.Equal(t, contained, aContained)
	total := unique + duplicates + contained

	assert.Equal(t, "//// Contained Summary ///////////////////", lines[currentLine+2])
	assertSummaryLine(t, "Total files:", total, lines[currentLine+4])
	assertSummaryLine(t, "Contained files:", contained, lines[currentLine+5])
	assertSummaryLine(t, "Duplicate files:", duplicates, lines[currentLine+6])
}

func AssertLogBlocks(t *testing.T, lines []string, blocks []LogBlock) (int, int, int) {
	currentLine := 0
	duplicates := 0
	contained := 0
	for _, block := range blocks {
		assert.Equal(t, block.firstLine(), lines[currentLine])
		currentLine++
		for _, relativePath := range block.relativePaths {
			assert.Equal(t, relativePath, lines[currentLine])
			currentLine++
		}
		switch block.category {
		case ilog.Duplicate:
			duplicates += len(block.relativePaths) - 1
		case ilog.Contained:
			contained += len(block.relativePaths)
		}
		currentLine++
	}
	return currentLine, duplicates, contained
}

func AssertLogFileNotExists(t *testing.T, dir string) {
	absoluteDir := filepath.Join(dir, integrity)
	files, err := os.ReadDir(absoluteDir)
	if err != nil {
		log.Fatal("could not read log dir", err)
	}
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".log" {
			assert.Fail(t, "log file should not exist", f.Name())
		}
	}
}

func assertSummaryLine(t *testing.T, expectedLabel string, expectedValue int, actualLine string) {
	assert.Len(t, actualLine, 42)
	values := strings.Split(actualLine, ":")
	assert.Len(t, values, 2)
	actualLabel := values[0] + ":"
	actualValue := strings.TrimSpace(values[1])
	assert.Equal(t, expectedLabel, actualLabel)
	assert.Equal(t, strconv.Itoa(expectedValue), actualValue, expectedLabel)
}
