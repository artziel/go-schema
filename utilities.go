package schema

import "strings"

type pair struct {
	Name  string
	Value string
}

func toPair(value string) pair {
	result := pair{}

	i := strings.Index(value, ":")
	if i > -1 && i < len([]rune(value)) {
		result.Name = value[:i]
		result.Value = strings.Trim(value[i+1:], " \n\r\t")
	} else {
		result.Name = value
	}

	result.Name = strings.Trim(strings.ToLower(result.Name), " \n\r\t")

	return result
}

func splitAt(s string, separator byte, quote byte) []string {
	res := []string{}
	var beg int
	var inString bool

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
