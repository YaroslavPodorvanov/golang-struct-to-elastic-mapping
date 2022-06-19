package generator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/YaroslavPodorvanov/golang-struct-to-elastic-mapping/converter"
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
			field              = value.Type().Field(i)
			propertyName, _, _ = strings.Cut(field.Tag.Get("json"), ",")
			propertyType       = g.kindConverter.Get(field.Type.Kind())
		)

		if propertyName == "" {
			propertyName = field.Name
		}

		if propertyType == "" {
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
