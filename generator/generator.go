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
	properties, err := g.properties(i)

	if err != nil {
		return nil, err
	}

	return wrap(properties)
}

func (g *Generator) properties(i any) ([]byte, error) {
	var (
		result = make([]byte, 0, 1024)

		value  = reflect.ValueOf(i).Elem()
		length = value.NumField()
	)

	result = append(result, '{')

	for i := 0; i < length; i++ {
		var (
			field                  = value.Type().Field(i)
			jsonPropertyName, _, _ = strings.Cut(field.Tag.Get("json"), ",")
			tags                   = parseTags(field.Tag.Get("es"))
			propertyType           = tags.typeName
		)

		if jsonPropertyName == "" {
			jsonPropertyName = field.Name
		}

		if propertyType == "" {
			propertyType = g.kindConverter.Get(field.Type.Kind())
		}

		if propertyType == "" {
			return nil, errors.New(fmt.Sprintf("cannot found %s", field.Type.String()))
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
