package common

import (
	"testing"

	"github.com/aicirt2012/fileintegrity/src/store/ilog"
	"github.com/stretchr/testify/assert"
)

func TestLogStyleMap(t *testing.T) {
	log1 := ilog.StyleLog{RelativePath: "b/c"}
	log2 := ilog.StyleLog{RelativePath: "a/b"}

	m := LogStyleMap{}
	m.Put(log1)
	m.Put(log2)

	expected := []ilog.StyleLog{log2, log1}
	assert.Len(t, m, 2)
	assert.Equal(t, expected, m.SortValues())
}
