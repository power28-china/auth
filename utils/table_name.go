package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gertd/go-pluralize"
)

// ConvertToTableName convert the type name of a object to a table name. lower case and plural format
func ConvertToTableName(obj interface{}) string {
	pluralize := pluralize.NewClient()

	tab := fmt.Sprintf("%s", reflect.TypeOf(obj))
	index := strings.Index(tab, ".") + 1
	return pluralize.Plural(strings.ToLower(tab[index:]))
}

func isUpper(b byte) bool {
	return 'A' <= b && b <= 'Z'
}

func isLower(b byte) bool {
	return 'a' <= b && b <= 'z'
}

func isDigit(b byte) bool {
	return '0' <= b && b <= '9'
}

func toLower(b byte) byte {
	if isUpper(b) {
		return b - 'A' + 'a'
	}
	return b
}

// CamelCaseToSnakeCase returns a snake case string from a camel case string.
func CamelCaseToSnakeCase(name string) string {
	var buf strings.Builder
	buf.Grow(len(name) * 2)

	for i := 0; i < len(name); i++ {
		buf.WriteByte(toLower(name[i]))
		if i != len(name)-1 && isUpper(name[i+1]) &&
			(isLower(name[i]) || isDigit(name[i]) ||
				(i != len(name)-2 && isLower(name[i+2]))) {
			buf.WriteByte('_')
		}
	}

	return buf.String()
}

// ConvertToMySqlTableName convert the type name of a object to a table name. lower case
func ConvertToMySqlTableName(obj interface{}) string {

	tab := fmt.Sprintf("%s", reflect.TypeOf(obj))
	index := strings.Index(tab, ".") + 1
	return fmt.Sprintf("%s", strings.ToLower(tab[index:]))
}
