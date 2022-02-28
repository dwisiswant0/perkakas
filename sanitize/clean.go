package sanitize

import "reflect"

// Clean implements deep-clean and/or sanitization of interface
func Clean(i interface{}, strict bool) interface{} {
	deepClean(reflect.ValueOf(i), setPolicy(strict))

	return i
}
