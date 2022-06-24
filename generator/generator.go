package generator

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/YaroslavPodorvanov/golang-struct-to-elastic-mapping/converter"
	"github.com/YaroslavPodorvanov/golang-struct-to-elastic-mapping/mapping"
)

type Generator struct {
	kindConverter *converter.KindConverter
}

func NewGenerator(kindConverter *converter.KindConverter) *Generator {
	return &Generator{kindConverter: kindConverter}
}

var DefaultGenerator = NewGenerator(converter.DefaultKindConverter())

func Generate(i any) ([]byte, error) {
	return DefaultGenerator.Generate(i)
}

func (g *Generator) Generate(i any) ([]byte, error) {
	properties, err := g.properties(reflect.ValueOf(i).Type().Elem())

	if err != nil {
		return nil, err
	}

	return wrap(properties)
}

func (g *Generator) properties(t reflect.Type) ([]byte, error) {
	var (
		result = make([]byte, 0, 1024)

		length = t.NumField()
	)

	result = append(result, '{')

MAIN:
	for i := 0; i < length; i++ {
		var (
			field                  = t.Field(i)
			fieldType              = field.Type
			jsonPropertyName, _, _ = strings.Cut(field.Tag.Get("json"), ",")
			tags, skip             = parseTags(field.Tag.Get("es"))
			propertyType           = tags.typeName
		)

		// @TODO: if last then incorrect JSON
		if skip {
			continue
		}

		if jsonPropertyName == "" {
			jsonPropertyName = field.Name
		}

		if propertyType == "" {
			propertyType = g.kindConverter.Get(fieldType.Kind())
		}

		if propertyType == "" {
			switch fieldType.Kind() {
			case reflect.Struct:
				var properties, err = g.properties(fieldType)
				if err != nil {
					return nil, err
				}

				var property, marshalErr = json.Marshal(&mapping.NestedProperties{
					Type:       "nested",
					Properties: properties,
				})
				if marshalErr != nil {
					return nil, marshalErr
				}

				result = append(result, fmt.Sprintf(`"%s":`, jsonPropertyName)...)
				result = append(result, property...)

				if !lastIndex(length, i) {
					result = append(result, ',')
				}

				continue MAIN
			case reflect.Slice:
				var properties, err = g.properties(fieldType.Elem())
				if err != nil {
					return nil, err
				}

				var property, marshalErr = json.Marshal(&mapping.NestedProperties{
					Type:       "nested",
					Properties: properties,
				})
				if marshalErr != nil {
					return nil, marshalErr
				}

				result = append(result, fmt.Sprintf(`"%s":`, jsonPropertyName)...)
				result = append(result, property...)

				if !lastIndex(length, i) {
					result = append(result, ',')
				}

				continue MAIN
			}
		}

		if propertyType == "" {
			return nil, errors.New(fmt.Sprintf("cannot found %s", fieldType.String()))
		}

		var property, err = json.Marshal(&mapping.Property{
			Type:  propertyType,
			Index: tags.index,
		})
		if err != nil {
			return nil, err
		}

		result = append(result, fmt.Sprintf(`"%s":`, jsonPropertyName)...)
		result = append(result, property...)

		if !lastIndex(length, i) {
			result = append(result, ',')
		}
	}

	result = append(result, '}')

	return result, nil
}
