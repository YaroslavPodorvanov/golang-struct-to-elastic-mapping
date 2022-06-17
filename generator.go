package mapping

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

func Generate(i any) ([]byte, error) {
	properties, err := generateBytes(i)

	if err != nil {
		return nil, err
	}

	return wrap(properties)
}

func generateBytes(i any) ([]byte, error) {
	var (
		result = make([]byte, 0, 1024)

		value  = reflect.ValueOf(i).Elem()
		length = value.NumField()
	)

	result = append(result, '{')

	for i := 0; i < length; i++ {
		var (
			field              = value.Type().Field(i)
			propertyName, _, _ = strings.Cut(field.Tag.Get("json"), ",")
			propertyType       string
		)

		if propertyName == "" {
			propertyName = field.Name
		}

		//	@TODO switch field.Type.Kind()
		//		1. keyword
		//		2. binary

		switch field.Type.Kind() {
		//	https://www.elastic.co/guide/en/elasticsearch/reference/current/number.html#number
		case reflect.Int64:
			propertyType = "long"
		case reflect.Int, reflect.Int32, reflect.Uint16:
			//		uint16 to integer
			propertyType = "integer"
		case reflect.Int16, reflect.Uint8:
			//		uint8 to short
			propertyType = "short"
		case reflect.Int8:
			//	from documentation:
			//		byte — A signed 8-bit integer with a minimum value of -128 and a maximum value of 127.
			propertyType = "byte"
		case reflect.Uint, reflect.Uint32, reflect.Uint64:
			propertyType = "unsigned_long"
		case reflect.Float32:
			propertyType = "float"
		case reflect.Float64:
			propertyType = "double"
		case reflect.Bool:
			propertyType = "boolean"
		case reflect.String:
			propertyType = "text"
		case reflect.TypeOf(time.Time{}).Kind():
			propertyType = "date"
		case reflect.TypeOf(map[string]interface{}{}).Kind():
			propertyType = "object"
		default:
			return nil, errors.New(fmt.Sprintf("cannot found %s", field.Type.String()))
		}

		result = append(result, fmt.Sprintf(`"%s":{"type":"%s","index":false}`, propertyName, propertyType)...)

		if !lastIndex(length, i) {
			result = append(result, ',')
		}
	}

	result = append(result, '}')

	return result, nil
}

func wrap(properties []byte) ([]byte, error) {
	return json.MarshalIndent(&Index{
		Mappings: IndexMappings{
			Properties: properties,
		},
	}, "", "  ")
}

func lastIndex(length, index int) bool {
	return length-1 == index
}
