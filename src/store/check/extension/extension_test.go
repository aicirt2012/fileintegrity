package extension

import (
	"testing"

	"github.com/aicirt2012/fileintegrity/src/store/file"
	"github.com/stretchr/testify/assert"
)

func TestCalcExtensionMap(t *testing.T) {
	input := []file.FileHash{
		{
			Size:         10,
			RelativePath: "/dir/a.jpg",
		},
		{
			Size:         30,
			RelativePath: "/dir/b.JPG",
		},
		{
			Size:         10,
			RelativePath: "/dir/b.png",
		},
	}
	expectedMap := extMap{
		".jpg": {
			name:  ".jpg",
			bytes: 40,
		},
		".png": {
			name:  ".png",
			bytes: 10,
		},
	}
	actualMap, actualTotalBytes := calcExtensionMap(input)
	assert.Equal(t, expectedMap, actualMap)
	assert.Equal(t, int64(50), actualTotalBytes)
}
