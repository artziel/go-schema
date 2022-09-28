package schema

import (
	"strconv"
	"strings"
)

type Tag struct {
	Values map[string]string
}

func (t *Tag) IsEmpty() bool {
	return len(t.Values) == 0
}

func (t *Tag) Exists(name string) bool {
	if _, exists := t.Values[strings.ToLower(name)]; exists {
		return true
	}
	return false
}

func (t *Tag) HasValue(name string) bool {
	if t.Exists(name) {
		return t.Values[strings.ToLower(name)] != ""
	}
	return false
}

func (t *Tag) GetString(name string) string {
	if t.Exists(name) {
		return t.Values[strings.ToLower(name)]
	}
	return ""
}

func (t *Tag) GetUint(name string) uint64 {
	v := t.GetString(name)

	if v != "" {
		value, _ := strconv.ParseUint(v, 10, 64)
		return value
	}

	return 0
}

func (t *Tag) GetInt(name string) int64 {
	v := t.GetString(name)

	if v != "" {
		value, _ := strconv.Atoi(v)
		return int64(value)
	}

	return 0
}

func (t *Tag) GetFloat(name string) float64 {
	v := t.GetString(name)

	if v != "" {
		value, _ := strconv.ParseFloat(v, 64)
		return value
	}

	return 0
}

func (t *Tag) GetSliceString(name string, separator byte) []string {
	result := []string{}
	v := t.GetString(name)
	if v == "" {
		return result
	}

	for _, itm := range splitAt(v, separator, '\'') {
		result = append(result, strings.Trim(itm, "'"))
	}

	return result
}

func (t *Tag) GetSliceInt(name string, separator byte) []int64 {
	values := t.GetSliceString(name, separator)
	if len(values) == 0 {
		return []int64{}
	}
	result := []int64{}

	for _, itm := range values {
		v, _ := strconv.Atoi(itm)
		result = append(result, int64(v))
	}

	return result
}

func (t *Tag) GetSliceFloat(name string, separator byte) []float64 {
	values := t.GetSliceString(name, separator)
	if len(values) == 0 {
		return []float64{}
	}
	result := []float64{}

	for _, itm := range values {
		v, _ := strconv.ParseFloat(itm, 64)
		result = append(result, v)
	}

	return result
}

func (t *Tag) GetSliceUint(name string, separator byte) []uint64 {
	values := t.GetSliceString(name, separator)
	if len(values) == 0 {
		return []uint64{}
	}
	result := []uint64{}

	for _, itm := range values {
		v, _ := strconv.ParseUint(itm, 10, 64)
		result = append(result, v)
	}

	return result
}

func ParseTag(value string) Tag {

	tag := Tag{
		Values: map[string]string{},
	}

	values := splitAt(value, ',', '\'')

	for _, v := range values {
		p := toPair(v)
		tag.Values[p.Name] = p.Value
	}

	return tag
}
