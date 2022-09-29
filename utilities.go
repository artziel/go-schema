package schema

import "strings"

/*
Structure used as return value by the function toPair
*/
type pair struct {
	Name  string
	Value string
}

/*
Split an string in 2 an return a pair struct value

The first part resulting from the split will be assigned to Name and
the second one to Value
*/
func toPair(value string, separator byte) pair {
	result := pair{}

	i := strings.Index(value, string(separator))
	if i > -1 && i < len([]rune(value)) {
		result.Name = value[:i]
		result.Value = strings.Trim(value[i+1:], " \n\r\t")
	} else {
		result.Name = value
	}

	result.Name = strings.Trim(strings.ToLower(result.Name), " \n\r\t")

	return result
}

/*
Split an string on N parts

The "separator" param indicates the character from which the string will be split, the separator
character will be ignored during the evaluation if it is between two characters equal to "quote"
parameter
*/
func splitAt(s string, separator byte, quote byte) []string {

	res := []string{}
	var beg int
	var inString bool

	s = strings.Trim(s, " \n\r\t")
	if len(s) == 0 {
		return res
	}

	for i := 0; i < len(s); i++ {
		if s[i] == separator && !inString {
			res = append(res, s[beg:i])
			beg = i + 1
		} else if s[i] == quote {
			if !inString {
				inString = true
			} else if i > 0 && s[i-1] != '\\' {
				inString = false
			}
		}
	}
	return append(res, s[beg:])
}
