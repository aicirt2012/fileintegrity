package extension

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderedValues(t *testing.T) {
	input := extMap{
		".jpg": {
			name:  ".jpg",
			bytes: 500,
		},
		".png": {
			name:  ".png",
			bytes: 499,
		},
		".gif": {
			name:  ".gif",
			bytes: 498,
		},
		".zzz": {
			name:  ".zzz",
			bytes: 498,
		},
	}
	expected := []ext{
		{
			name:  ".jpg",
			bytes: 500,
		},
		{
			name:  ".png",
			bytes: 499,
		},
		{
			name:  ".gif",
			bytes: 498,
		},
		{
			name:  ".zzz",
			bytes: 498,
		},
	}
	actual := input.orderedValues()
	assert.Equal(t, expected, actual)
}
