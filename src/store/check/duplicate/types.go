package duplicate

import (
	"sort"
	"strconv"

	"github.com/aicirt2012/fileintegrity/src/store/file"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type UniqueMap map[string]file.FileHashs

func (m UniqueMap) Has(key string) bool {
	if _, ok := m[key]; ok {
		return true
	}
	return false
}

func (m UniqueMap) add(fh file.FileHash) {
	m[key(fh)] = append(m[key(fh)], fh)
}

func (m *UniqueMap) Remove(key string) {
	delete(*m, key)
}

func (m UniqueMap) orderedValues() []file.FileHashs {
	result := []file.FileHashs{}
	keys := maps.Keys(m)
	slices.Sort(keys)
	for _, key := range keys {
		if fhs, ok := m[key]; ok {
			sort.Sort(fhs)
			result = append(result, fhs)
		}
	}
	return result
}

func key(fh file.FileHash) string {
	return fh.Hash + strconv.FormatInt(fh.Size, 10)
}
