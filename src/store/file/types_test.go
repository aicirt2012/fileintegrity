package file

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefragmentedMap(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name     string
		in       FileHashs
		expected FileHashs
	}{
		{
			name:     "Unchanged empty",
			in:       FileHashs{},
			expected: FileHashs{},
		},
		{
			name: "Unchanged single entry",
			in: FileHashs{
				FileHash{
					Hash:         "a2bd6fecf7366e8c9722371b4c06f78adfbb46a92e483b42324e4b4d27b5232e",
					Created:      now,
					ModTime:      now.Add(-time.Hour),
					Size:         100,
					RelativePath: `basedir\subdir\filename.jpg`,
				},
			},
			expected: FileHashs{
				FileHash{
					Hash:         "a2bd6fecf7366e8c9722371b4c06f78adfbb46a92e483b42324e4b4d27b5232e",
					Created:      now,
					ModTime:      now.Add(-time.Hour),
					Size:         100,
					RelativePath: `basedir\subdir\filename.jpg`,
				},
			},
		},
		{
			name: "Update one entry",
			in: FileHashs{
				FileHash{
					Hash:         "a2bd6fecf7366e8c9722371b4c06f78adfbb46a92e483b42324e4b4d27b5232e",
					Created:      now.Add(-time.Hour),
					ModTime:      now.Add(-time.Hour),
					Size:         100,
					RelativePath: `basedir\subdir\filename.jpg`,
				},
				FileHash{
					Hash:         "16f3bb4b94a001e4de4fddbef555cf6cf02022b4f8180023bc8b905fb3e7e373",
					Created:      now,
					ModTime:      now,
					Size:         120,
					RelativePath: `basedir\subdir\filename.jpg`,
				},
				FileHash{
					Hash:         "02914e8ea7ffb5632bbb84fe1453372aaa437d196ba7da5505f4eaed0cfc10f2",
					Created:      now,
					ModTime:      now.Add(-2 * time.Hour),
					Size:         300,
					RelativePath: `basedir\subdir\noisefile.jpg`,
				},
			},
			expected: FileHashs{
				FileHash{
					Hash:         "16f3bb4b94a001e4de4fddbef555cf6cf02022b4f8180023bc8b905fb3e7e373",
					Created:      now,
					ModTime:      now,
					Size:         120,
					RelativePath: `basedir\subdir\filename.jpg`,
				},
				FileHash{
					Hash:         "02914e8ea7ffb5632bbb84fe1453372aaa437d196ba7da5505f4eaed0cfc10f2",
					Created:      now,
					ModTime:      now.Add(-2 * time.Hour),
					Size:         300,
					RelativePath: `basedir\subdir\noisefile.jpg`,
				},
			},
		},
		{
			name: "Delete one entry",
			in: FileHashs{
				FileHash{
					Hash:         "a2bd6fecf7366e8c9722371b4c06f78adfbb46a92e483b42324e4b4d27b5232e",
					Created:      now,
					ModTime:      now.Add(-time.Hour),
					Size:         100,
					RelativePath: `basedir\subdir\filename.jpg`,
				},
				FileHash{
					Hash:         "16f3bb4b94a001e4de4fddbef555cf6cf02022b4f8180023bc8b905fb3e7e373",
					Created:      now.Add(-time.Hour),
					ModTime:      now.Add(-2 * time.Hour),
					Size:         100,
					RelativePath: `basedir\subdir\todelete.jpg`,
				},
				FileHash{
					Hash:         "0000000000000000000000000000000000000000000000000000000000000000",
					Created:      now,
					ModTime:      now.Add(-2 * time.Hour),
					Size:         100,
					RelativePath: `basedir\subdir\todelete.jpg`,
				},
			},
			expected: FileHashs{
				FileHash{
					Hash:         "a2bd6fecf7366e8c9722371b4c06f78adfbb46a92e483b42324e4b4d27b5232e",
					Created:      now,
					ModTime:      now.Add(-time.Hour),
					Size:         100,
					RelativePath: `basedir\subdir\filename.jpg`,
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.in.DefragmentedMap()
			assert.Equal(t, len(c.expected), len(actual))
			for _, fh := range c.expected {
				actualFileHash, exists := actual[fh.RelativePath]
				assert.True(t, exists)
				assert.Equal(t, fh, actualFileHash)
			}
		})
	}
}
