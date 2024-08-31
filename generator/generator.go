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
		length         = t.NumField()
		rootProperties = make([][]byte, 0, length)
	)

MAIN:
	for i := 0; i < length; i++ {
		var (
			field                  = t.Field(i)
			fieldType              = field.Type
			jsonPropertyName, _, _ = strings.Cut(field.Tag.Get("json"), ",")
			tags, skip             = parseTags(field.Tag.Get("es"))
			propertyType           = tags.typeName
		)

		if skip {
			continue MAIN
		}

		if jsonPropertyName == "" {
			jsonPropertyName = field.Name
		}

		if propertyType == "" {
			propertyType = g.kindConverter.Get(fieldType.Kind())
		}

		if propertyType == "" {
			switch fieldType.Kind() {
			case reflect.Pointer:
				properties, err := g.properties(fieldType.Elem())
				if err != nil {
					return nil, err
				}

				property, err := g.toProperty(jsonPropertyName, &mapping.NestedProperties{
					Type:       "nested",
					Properties: properties,
				})
				if err != nil {
					return nil, err
				}

				rootProperties = append(rootProperties, property)

				continue MAIN
			case reflect.Struct:
				properties, err := g.properties(fieldType)
				if err != nil {
					return nil, err
				}

				property, err := g.toProperty(jsonPropertyName, &mapping.NestedProperties{
					Type:       "nested",
					Properties: properties,
				})
				if err != nil {
					return nil, err
				}

				rootProperties = append(rootProperties, property)

				continue MAIN
			case reflect.Slice:
				if fieldType.Elem().Kind() == reflect.Struct {
					properties, err := g.properties(fieldType.Elem())
					if err != nil {
						return nil, err
					}

					property, err := g.toProperty(jsonPropertyName, &mapping.NestedProperties{
						Type:       "nested",
						Properties: properties,
					})
					if err != nil {
						return nil, err
					}

					rootProperties = append(rootProperties, property)

					continue MAIN
				}

				propertyType = g.kindConverter.Get(fieldType.Elem().Kind())
			}
		}

		if propertyType == "" {
			return nil, errors.New(fmt.Sprintf("cannot found %s", fieldType.String()))
		}

		property, err := g.toProperty(jsonPropertyName, &mapping.Property{
			Type:  propertyType,
			Index: tags.index,
		})
		if err != nil {
			return nil, err
		}

		rootProperties = append(rootProperties, property)
	}

	return g.result(rootProperties)
}

func (g *Generator) result(properties [][]byte) ([]byte, error) {
	var (
		length = len(properties)
		result = make([]byte, 0, g.estimateCapacity(properties))
	)

	result = append(result, '{')

	for i, property := range properties {
		result = append(result, property...)

		if i < length-1 {
			result = append(result, ',')
		}
	}

	result = append(result, '}')

	return result, nil
}

func (g *Generator) estimateCapacity(properties [][]byte) int {
	var (
		result int
	)

	result += 2 // for '{' and '}'

	result += len(properties) - 1 // for commas

	for _, property := range properties {
		result += len(property)
	}

	return result
}

func (g *Generator) toProperty(name string, property any) ([]byte, error) {
	content, err := json.Marshal(property)
	if err != nil {
		return nil, err
	}

	var result []byte

	result = append(result, fmt.Sprintf(`%q:`, name)...)
	result = append(result, content...)

	return result, nil
}
