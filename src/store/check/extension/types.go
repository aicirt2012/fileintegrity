package extension

import (
	"sort"

	"golang.org/x/exp/maps"
)

type ext struct {
	name       string
	bytes      int64
	percentage float64
}

type extMap map[string]*ext

func (m extMap) orderedValues() []ext {
	values := []ext{}
	for _, value := range maps.Values(m) {
		values = append(values, *value)
	}
	sort.Slice(values, func(i, j int) bool {
		if values[i].bytes == values[j].bytes {
			return values[i].name < values[j].name
		}
		return values[i].bytes > values[j].bytes
	})
	return values
}
