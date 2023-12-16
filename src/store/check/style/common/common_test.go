package common

import (
	"testing"

	"github.com/aicirt2012/fileintegrity/tests/common"
	"github.com/stretchr/testify/assert"
)

func TestSplitPath(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Directory path",
			input:    "root/sub",
			expected: []string{"root", "sub"},
		},
		{
			name:     "File Path",
			input:    "root/sub/f.txt",
			expected: []string{"root", "sub", "f.txt"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := SplitPath(common.NormalizePath(c.input))
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestSplitDirs(t *testing.T) {
	input := common.NormalizePath("root/sub/f.txt")
	actual := SplitDirs(input)
	expected := []string{"root", "sub"}
	assert.Equal(t, expected, actual)
}

func TestMinimalPath(t *testing.T) {
	cases := []struct {
		name     string
		i        int
		sections []string
		expected string
	}{
		{
			name:     "Empty path",
			i:        0,
			sections: []string{"root", "sub", "x"},
			expected: "",
		},
		{
			name:     "Partial path",
			i:        2,
			sections: []string{"root", "sub", "x"},
			expected: common.NormalizePath("root/sub"),
		},
		{
			name:     "Full part",
			i:        2,
			sections: []string{"root", "sub"},
			expected: common.NormalizePath("root/sub"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := MinimalPath(c.i, c.sections)
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestJoinWithAnd(t *testing.T) {
	cases := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "Empty case",
			input:    []string{},
			expected: "",
		},
		{
			name:     "One item case",
			input:    []string{"a"},
			expected: "a",
		},
		{
			name:     "Two items case",
			input:    []string{"a", "b"},
			expected: "a and b",
		},
		{
			name:     "Three items case",
			input:    []string{"a", "b", "c"},
			expected: "a, b and c",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := JoinWithAnd(c.input)
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestQuote(t *testing.T) {
	input := []string{"a", "b"}
	expected := []string{"'a'", "'b'"}
	actual := Quote(input)
	assert.Equal(t, expected, actual)
}

func TestPluralS(t *testing.T) {
	assert.Equal(t, "", PluralS([]string{}))
	assert.Equal(t, "", PluralS([]string{"a"}))
	assert.Equal(t, "s", PluralS([]string{"a", "b"}))
}
