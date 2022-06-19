package converter

import (
	"reflect"
	"time"
)

type KindConverter struct {
	kindMap map[reflect.Kind]string
}

func NewKindConverter(kindMap map[reflect.Kind]string) *KindConverter {
	return &KindConverter{kindMap: kindMap}
}

func (k *KindConverter) Set(key reflect.Kind, value string) {
	k.kindMap[key] = value
}

func (k *KindConverter) Get(key reflect.Kind) string {
	return k.kindMap[key]
}

func (k *KindConverter) Clone() *KindConverter {
	var cloneKindMap = make(map[reflect.Kind]string, len(k.kindMap))

	for key, value := range k.kindMap {
		cloneKindMap[key] = value
	}

	return NewKindConverter(cloneKindMap)
}

func DefaultKindConverter() *KindConverter {
	return defaultKindConverter.Clone()
}

// https://www.elastic.co/guide/en/elasticsearch/reference/current/number.html#number
var defaultKindConverter = NewKindConverter(map[reflect.Kind]string{
	reflect.Bool:                       "boolean",
	reflect.Int:                        "integer",
	reflect.Int8:                       "byte", // A signed 8-bit integer with a minimum value of -128 and a maximum value of 127.
	reflect.Int16:                      "short",
	reflect.Int32:                      "integer",
	reflect.Int64:                      "long",
	reflect.Uint:                       "unsigned_long",
	reflect.Uint8:                      "short",
	reflect.Uint16:                     "integer",
	reflect.Uint32:                     "unsigned_long",
	reflect.Uint64:                     "unsigned_long",
	reflect.Float32:                    "float",
	reflect.Float64:                    "double",
	reflect.String:                     "text",
	reflect.TypeOf(time.Time{}).Kind(): "date",
})
