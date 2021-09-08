package graphql

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	graphqlgo "github.com/graphql-go/graphql"
	"github.com/serenize/snaker"
)

const (
	formatError = `[ResolveNullable]`
)

func getFieldStruct(fieldName string) (result string) {
	result = strings.Title(fieldName)
	if gqlInst.IsSnakeCase {
		result = snaker.SnakeToCamel(strings.ToLower(fieldName))
	}
	return
}

func resolveFieldValue(fieldValue interface{}) (result interface{}) {
	switch value := fieldValue.(type) {
	case sql.NullString:
		if value.Valid {
			result = value.String
		}
	case *sql.NullString:
		if value != nil && value.Valid {
			result = value.String
		}
	case sql.NullInt32:
		if value.Valid {
			result = value.Int32
		}
	case *sql.NullInt32:
		if value != nil && value.Valid {
			result = value.Int32
		}
	case sql.NullInt64:
		if value.Valid {
			result = value.Int64
		}
	case *sql.NullInt64:
		if value != nil && value.Valid {
			result = value.Int64
		}
	case sql.NullFloat64:
		if value.Valid {
			result = value.Float64
		}
	case *sql.NullFloat64:
		if value != nil && value.Valid {
			result = value.Float64
		}
	case sql.NullBool:
		if value.Valid {
			result = value.Bool
		}
	case *sql.NullBool:
		if value != nil && value.Valid {
			result = value.Bool
		}
	case sql.NullTime:
		if value.Valid {
			result = value.Time
		}
	case *sql.NullTime:
		if value != nil && value.Valid {
			result = value.Time
		}
	case *string:
		if value != nil {
			result = *value
		}
	case []byte:
		result = string(value)
	case *[]byte:
		if value != nil {
			result = string(*value)
		}
	default:
		if fieldValue != nil {
			result = fieldValue
		}
	}
	return
}

// ResolveNullable resolve the null field for any type in the struct
func ResolveNullable(p graphqlgo.ResolveParams) (result interface{}, err error) {
	var (
		fieldName   string
		fieldStruct string
		val         reflect.Value
		field       reflect.Value
		errMsg      string
	)
	fieldName = p.Info.FieldName
	val = reflect.ValueOf(p.Source)
	if !val.IsValid() {
		errMsg = fmt.Sprintf("%s Failed when try to get value: %s", formatError, fieldName)
		err = errors.New(errMsg)
		return
	}

	fieldStruct = getFieldStruct(fieldName)
	field = val.FieldByName(fieldStruct)
	if !field.IsValid() {
		errMsg = fmt.Sprintf("%s Missing field: %s", formatError, fieldStruct)
		err = errors.New(errMsg)
		return
	}

	result = resolveFieldValue(field.Interface())

	return
}
