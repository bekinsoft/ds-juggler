/*
 * @author    Emmanuel Kofi Bessah
 * @email     ekbessah@uew.edu.gh
 * @created   Sat Jul 07 2018 23:00:35
 */

package juggler

import (
	"regexp"
	"strings"
)

// JoinStrings ... Joins an array of string
func JoinStrings(values ...string) string {
	return strings.Join(values, "")
}

// ToSnake ... Converts a string to snake_case
// func ToSnake(s string) string {
// 	s = addWordBoundariesToNumbers(s)
// 	s = strings.Trim(s, " ")
// 	n := ""
// 	for i, v := range s {
// 		// treat acronyms as words, eg for JSONData -> JSON is a whole word
// 		nextIsCapital := false
// 		if i+1 < len(s) {
// 			w := s[i+1]
// 			nextIsCapital = w >= 'A' && w <= 'Z'
// 		}
// 		if i > 0 && v >= 'A' && v <= 'Z' && n[len(n)-1] != '_' && !nextIsCapital {
// 			// add underscore if next letter is a capital
// 			n += "_" + string(v)
// 		} else if v == ' ' {
// 			// replace spaces with underscores
// 			n += "_"
// 		} else {
// 			n = n + string(v)
// 		}
// 	}
// 	n = strings.ToLower(n)
// 	return n
// }

var numberSequence = regexp.MustCompile(`([a-zA-Z])(\d+)([a-zA-Z]?)`)
var numberReplacement = []byte(`$1 $2 $3`)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func addWordBoundariesToNumbers(s string) string {
	b := []byte(s)
	b = numberSequence.ReplaceAll(b, numberReplacement)
	return string(b)
}

// ToSnakeCase ...
func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// ToCamel ... Converts string to CamelCase
func ToCamel(s string) string {
	return toCamelInitCase(s, true)
}

// ToLowerCamel ... Converts string to camelCase
func ToLowerCamel(s string) string {
	return toCamelInitCase(s, false)
}

// Converts a string to CamelCase
func toCamelInitCase(s string, initCase bool) string {
	s = addWordBoundariesToNumbers(s)
	s = strings.Trim(s, " ")
	n := ""
	capNext := initCase
	for _, v := range s {
		if v >= 'A' && v <= 'Z' {
			n += string(v)
		}
		if v >= '0' && v <= '9' {
			n += string(v)
		}
		if v >= 'a' && v <= 'z' {
			if capNext {
				n += strings.ToUpper(string(v))
			} else {
				n += string(v)
			}
		}
		if v == '_' || v == ' ' || v == '-' {
			capNext = true
		} else {
			capNext = false
		}
	}
	return n
}

func stringListToLowerCase(arr []string) []string {
	s := []string{}
	for _, field := range arr {
		s = append(s, ToSnakeCase(field))
	}

	return s
}
