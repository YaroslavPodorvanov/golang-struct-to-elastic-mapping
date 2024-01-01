package main

import (
	"go/ast"
	go_parser "go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"golang.org/x/tools/go/ast/inspector"
)

type parser struct {
	parserType parserType
	path       string
}

type parserType string

const (
	parserTypeFile   = "file"
	parserTypeFolder = "folder"
)

type structMappingTask struct {
	structName string
	structType *ast.StructType
}

func (p *parser) do() (map[string]interface{}, error) {
	files, err := getFilesFromPath(p.parserType, p.path)
	if err != nil {
		return nil, err
	}

	var convertedStructs = make(map[string]interface{})
	for _, filePath := range files {
		astInFile, err := go_parser.ParseFile(token.NewFileSet(), filePath, nil, go_parser.ParseComments)
		if err != nil {
			return nil, err
		}
		i := inspector.New([]*ast.File{astInFile})
		iFilter := []ast.Node{
			&ast.GenDecl{},
		}

		var genTasks []structMappingTask
		i.Nodes(iFilter, func(n ast.Node, push bool) (proceed bool) {
			genDecl := n.(*ast.GenDecl)
			if genDecl.Doc == nil {
				return false
			}
			typeSpec, ok := genDecl.Specs[0].(*ast.TypeSpec)
			if !ok {
				return false
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				return false
			}
			for _, comment := range genDecl.Doc.List {
				switch comment.Text {
				case "//elasticgen:json":
					genTasks = append(genTasks, structMappingTask{
						structName: typeSpec.Name.Name,
						structType: structType,
					})
				}
			}
			return false
		})

		for _, task := range genTasks {
			newStruct := processStructFields(task.structType)
			instance := reflect.New(newStruct).Interface()
			convertedStructs[task.structName] = instance
		}
	}

	return convertedStructs, nil
}

func getFilesFromPath(parserType parserType, path string) ([]string, error) {
	var files []string

	switch parserType {
	case parserTypeFile:
		files = append(files, path)
	case parserTypeFolder:
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return files, nil
}

func processStructFields(t *ast.StructType) reflect.Type {
	var fields []reflect.StructField
	for _, field := range t.Fields.List {
		fieldType := convertAstToReflect(field.Type)
		_, fieldTags, _ := strings.Cut(field.Tag.Value, "`")
		structField := reflect.StructField{
			Name: field.Names[0].Name,
			Type: fieldType,
			Tag:  reflect.StructTag(fieldTags),
		}
		fields = append(fields, structField)
	}
	return reflect.StructOf(fields)
}

func convertAstToReflect(astType ast.Expr) reflect.Type {
	switch t := astType.(type) {
	case *ast.Ident:
		// Handle named types (e.g. structs)
		if t.Obj != nil && t.Obj.Decl != nil {
			if typeSpec, ok := t.Obj.Decl.(*ast.TypeSpec); ok {
				if structType, ok := typeSpec.Type.(*ast.StructType); ok {
					return processStructFields(structType)
				}
			}
		}
		if t, ok := defaultTypes[t.Name]; ok {
			return t
		}
		log.Printf("Unsupported field type: %s\n", reflect.TypeOf(t.Name).String())
		return nil
	case *ast.StructType:
		return processStructFields(t)
	case *ast.ArrayType:
		return reflect.SliceOf(convertAstToReflect(t.Elt))
	case *ast.MapType:
		return reflect.MapOf(convertAstToReflect(t.Key), convertAstToReflect(t.Value))
	case *ast.StarExpr:
		return convertAstToReflect(t.X)
	case *ast.SelectorExpr:
		return reflect.TypeOf(t.Sel.Name)
	case *ast.InterfaceType:
		return reflect.TypeOf(t)
	default:
		log.Printf("Unsupported type: %s\n", reflect.TypeOf(astType).String())
		return nil
	}
}

var defaultTypes = map[string]reflect.Type{
	"int":     reflect.TypeOf(0),
	"int8":    reflect.TypeOf(int8(0)),
	"int16":   reflect.TypeOf(int16(0)),
	"int32":   reflect.TypeOf(int32(0)),
	"int64":   reflect.TypeOf(int64(0)),
	"uint":    reflect.TypeOf(uint(0)),
	"uint8":   reflect.TypeOf(uint8(0)),
	"uint16":  reflect.TypeOf(uint16(0)),
	"uint32":  reflect.TypeOf(uint32(0)),
	"uint64":  reflect.TypeOf(uint64(0)),
	"float32": reflect.TypeOf(float32(0.0)),
	"float64": reflect.TypeOf(0.0),
	"string":  reflect.TypeOf(""),
	"bool":    reflect.TypeOf(false),
}
