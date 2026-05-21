package str

import (
	"encoding/json"
	"reflect"
	"strings"
)

func Trim(s string) string {
	return strings.TrimSpace(s)
}

func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

func Contains(s, sub string) bool {
	return strings.Contains(s, sub)
}

func ToJSONString(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

func ToString(v any) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(v.(string))
}

func Equal(a, b any) bool {
	return reflect.DeepEqual(a, b)
}

func IsEmptyString(s string) bool {
	return strings.TrimSpace(s) == ""
}
