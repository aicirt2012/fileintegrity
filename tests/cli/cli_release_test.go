//go:build release

package cli

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	actual := executeCli([]string{"-v"})
	actual = strings.TrimSpace(actual)
	expected := `^fileintegrity version \d+\.\d+\.\d+$`
	assert.Regexp(t, expected, actual, "version pattern invalid")
}

func TestLicense(t *testing.T) {
	actual := executeCli([]string{"license"})
	assert.Contains(t, actual, "Felix Michel", "license missing")
	assert.Contains(t, actual, "github.com/gocarina/gocsv", "transitive license missing")
	assert.True(t, len(actual) > 10000, "license text too short")
}
