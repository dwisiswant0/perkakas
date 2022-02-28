package sanitize

import (
	"reflect"

	"github.com/microcosm-cc/bluemonday"
)

var policy *bluemonday.Policy

func deepClean(v reflect.Value, p *bluemonday.Policy) reflect.Value {
	if v.Kind() != reflect.Ptr {
		return v
	}

	v = v.Elem()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)

		// Currently string & bytes are supported
		switch f.Kind() {
		case reflect.String:
			f.SetString(p.Sanitize(f.String()))
		case reflect.Slice:
			switch f.Type().Elem().Kind() {
			case reflect.Uint8: // []byte handler
				f.SetBytes(p.SanitizeBytes(f.Bytes()))
			case reflect.String: // []string handler
				for j := 0; j < f.Len(); j++ {
					s := f.Index(j)
					s.SetString(p.Sanitize(s.String()))
				}
			}
		}
	}

	return v
}

func setPolicy(strict bool) *bluemonday.Policy {
	if strict {
		policy = bluemonday.StrictPolicy()
	} else {
		policy = bluemonday.UGCPolicy()
	}

	return policy
}
