package schema

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
)

type fieldErrors []error

func (se fieldErrors) MarshalJSON() ([]byte, error) {
	data := []byte("[")
	for i, err := range se {
		if i != 0 {
			data = append(data, ',')
		}

		j, err := json.Marshal(err.Error())
		if err != nil {
			return nil, err
		}

		data = append(data, j...)
	}
	data = append(data, ']')

	return data, nil
}

type schemaTag struct {
	Name       string
	Required   bool
	Require    []string
	RestrictTo []string
	Regex      *regexp.Regexp
	MinLength  uint
	MaxLength  uint
	Min        float64
	Max        float64
}

type Field struct {
	Errors fieldErrors `json:"errors"`
}

type Result struct {
	Fields map[string]Field `json:"fields"`
}

func doubleQuotesSplitAt(s string, separator byte) []string {
	res := []string{}
	var beg int
	var inString bool

	for i := 0; i < len(s); i++ {
		if s[i] == separator && !inString {
			res = append(res, s[beg:i])
			beg = i + 1
		} else if s[i] == '"' {
			if !inString {
				inString = true
			} else if i > 0 && s[i-1] != '\\' {
				inString = false
			}
		}
	}
	return append(res, s[beg:])
}

func singleQuotesSplitAt(s string, separator byte) []string {
	res := []string{}
	var beg int
	var inString bool

	for i := 0; i < len(s); i++ {
		if s[i] == separator && !inString {
			res = append(res, s[beg:i])
			beg = i + 1
		} else if s[i] == '\'' {
			if !inString {
				inString = true
			} else if i > 0 && s[i-1] != '\\' {
				inString = false
			}
		}
	}
	return append(res, s[beg:])
}

func parse_tag(tag string) (schemaTag, error) {
	st := schemaTag{}

	regexName := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z_0-9]+$`)

	params := doubleQuotesSplitAt(tag, ',')
	for _, param := range params {
		pair := doubleQuotesSplitAt(param, ':')
		name := strings.ToLower(strings.Trim(pair[0], " \n\r\t"))
		value := ""
		if len(pair) > 1 {
			value = strings.Trim(pair[1], " \n\r\t")
		}
		switch name {
		case "required":
			st.Required = true
		case "restrictto":
			if value == "" {
				return st, ErrTagRestrictToMissingValue
			} else {
				st.RestrictTo = singleQuotesSplitAt(value, ' ')
			}
		case "require":
			if value == "" {
				return st, ErrTagRequireMissingValue
			} else {
				st.Require = doubleQuotesSplitAt(value, ' ')
			}
		case "name":
			if value == "" {
				return st, ErrTagNameMissingValue
			} else {
				if !regexName.MatchString(value) {
					return st, ErrTagNameValueInvalid
				}
				st.Name = value
			}
		case "regex":
			if value == "" {
				return st, ErrTagRegexMissingValue
			} else {
				r, err := regexp.Compile(value)
				if err != nil {
					return st, ErrTagRegexValueFailToCompile
				}
				st.Regex = r
			}
		case "minlength":
			if value == "" {
				return st, ErrTagMinLengthMissingValue
			} else {
				v, _ := strconv.Atoi(value)
				st.MinLength = uint(v)
			}
		case "maxlength":
			if value == "" {
				return st, ErrTagMinLengthMissingValue
			} else {
				v, _ := strconv.Atoi(value)
				st.MaxLength = uint(v)
			}
		case "min":
			if value == "" {
				return st, ErrTagMinMissingValue
			} else {
				v, _ := strconv.ParseFloat(value, 64)
				st.Min = float64(v)
			}
		case "max":
			if value == "" {
				return st, ErrTagMaxMissingValue
			} else {
				v, _ := strconv.ParseFloat(value, 64)
				st.Max = float64(v)
			}
		}
	}

	return st, nil
}
