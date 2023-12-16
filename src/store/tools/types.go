package tools

import (
	"sort"
	"strconv"

	"github.com/aicirt2012/fileintegrity/src/store/file"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type uniqueMap map[string]file.FileHashs

func (m uniqueMap) has(key string) bool {
	if _, ok := m[key]; ok {
		return true
	}
	return false
}

func (m uniqueMap) add(fh file.FileHash) {
	m[key(fh)] = append(m[key(fh)], fh)
}

func (m *uniqueMap) remove(key string) {
	delete(*m, key)
}

func (m uniqueMap) orderedValues() []file.FileHashs {
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
