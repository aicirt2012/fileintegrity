package ilog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTitle(t *testing.T) {
	actual := title(Upsert)
	expect := "\n\n//// Upsert Summary //////////////////////\n"
	assert.Equal(t, expect, actual)
}

func TestLine(t *testing.T) {
	actual := line("label", "%.2f", 3.456)
	expect := "label                                 3.46\n"
	assert.Equal(t, expect, actual)
}

func TestPadRight(t *testing.T) {
	assert.Equal(t, "l--", padRight("l", "-", 3))
	assert.Equal(t, "left", padRight("left", "-", 4))
	assert.Equal(t, "left", padRight("left", "-", -1))
}

func TestPadMiddle(t *testing.T) {
	assert.Equal(t, "r--l", padMiddle("r", "l", "-", 4))
	assert.Equal(t, "rile", padMiddle("ri", "le", "-", 4))
	assert.Equal(t, "rile", padMiddle("ri", "le", "-", -1))
}
