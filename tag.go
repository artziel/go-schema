package schema

import (
	"strconv"
	"strings"
)

/*
Structure that will contain the values of the parsed tag
*/
type Tag struct {
	Values map[string]string
}

/*
Check if the Tag structure contains values parsed or not
*/
func (t *Tag) IsEmpty() bool {
	return len(t.Values) == 0
}

/*
Check if the Tag structure contains the specified entry
*/
func (t *Tag) Exists(name string) bool {
	if _, exists := t.Values[strings.ToLower(name)]; exists {
		return true
	}
	return false
}

/*
Check if the Tag structure entry has a value
*/
func (t *Tag) HasValue(name string) bool {
	if t.Exists(name) {
		return t.Values[strings.ToLower(name)] != ""
	}
	return false
}

/*
Return the value of the specified entry as string
*/
func (t *Tag) GetString(name string) string {
	if t.Exists(name) {
		return t.Values[strings.ToLower(name)]
	}
	return ""
}

/*
Return the value of the specified entry as uint64
*/
func (t *Tag) GetUint(name string) uint64 {
	v := t.GetString(name)

	if v != "" {
		value, _ := strconv.ParseUint(v, 10, 64)
		return value
	}

	return 0
}

/*
Return the value of the specified entry as int64
*/
func (t *Tag) GetInt(name string) int64 {
	v := t.GetString(name)

	if v != "" {
		value, _ := strconv.Atoi(v)
		return int64(value)
	}

	return 0
}

/*
Return the value of the specified entry as float64
*/
func (t *Tag) GetFloat(name string) float64 {
	v := t.GetString(name)

	if v != "" {
		value, _ := strconv.ParseFloat(v, 64)
		return value
	}

	return 0
}

/*
Return an Slice of string of the specified entry. The separator parameter indicate
the character that will be used for the entry raw value string split
*/
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

/*
Return an Slice of int64 of the specified entry. The separator parameter indicate
the character that will be used for the entry raw value string split
*/
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

/*
Return an Slice of float64 of the specified entry. The separator parameter indicate
the character that will be used for the entry raw value string split
*/
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

/*
Return an Slice of uint64 of the specified entry. The separator parameter indicate
the character that will be used for the entry raw value string split
*/
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

/*
Parse an string and return a populated Tag structure value
*/
func ParseTag(value string) Tag {

	tag := Tag{
		Values: map[string]string{},
	}

	values := splitAt(value, ',', '\'')

	for _, v := range values {
		p := toPair(v, ':')
		tag.Values[p.Name] = p.Value
	}

	return tag
}
