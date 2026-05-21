package str

import (
	"reflect"
	"regexp"
	"strings"

	"golang.org/x/text/unicode/norm"
)

func Normalize(s string) string {
	s = strings.ToUpper(s)
	s = norm.NFD.String(s)

	re := regexp.MustCompile(`[^A-Z0-9 ]+`)
	s = re.ReplaceAllString(s, "")

	re2 := regexp.MustCompile(`\s+`)
	return re2.ReplaceAllString(strings.TrimSpace(s), " ")
}

func ToMap(obj any) map[string]any {
	result := map[string]any{}

	v := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)

	for i := 0; i < v.NumField(); i++ {
		result[t.Field(i).Name] = v.Field(i).Interface()
	}

	return result
}
