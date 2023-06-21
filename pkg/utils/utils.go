package utils

import (
	"reflect"
	"strings"
)

// IsEmpty will check for given data is empty as per the go documentation
func IsEmpty(val interface{}) bool {

	//First check normal definitions of empty
	if val == nil {
		return true
	} else if val == "" {
		return true
	} else if val == false {
		return true
	}

	reflectVal := reflect.ValueOf(val)

	switch reflectVal.Kind() {
	case reflect.Int:
		return val.(int) == 0

	case reflect.Int64:
		return val.(int64) == 0

	case reflect.String:
		return strings.TrimSpace(val.(string)) == ""

	case reflect.Map:
		fallthrough
	case reflect.Slice:
		return reflectVal.IsNil() || reflectVal.Len() == 0

	case reflect.Interface, reflect.Ptr:
		if reflectVal.IsNil() {
			return true
		}
		return IsEmpty(reflectVal.Elem().Interface())

	case reflect.Struct:
		copyStruct := reflect.New(reflect.TypeOf(val)).Elem().Interface()
		if reflect.DeepEqual(val, copyStruct) {
			return true
		}
	}

	return false
}
