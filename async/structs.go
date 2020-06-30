// Copyright 2020 Torben Schinke
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package async

import (
	"fmt"
	"github.com/golangee/openapi-client/internal/gen"
	v3 "github.com/golangee/openapi/v3"
)

func emitTypes(f *gen.GoGenFile, doc *v3.Document) error {
	if doc.Components == nil {
		return nil
	}

	for _, name := range gen.SortedKeys(doc.Components.Schemas) {
		if name == "Error" {
			continue // we ignore our own build-in type
		}
		schema := doc.Components.Schemas[name]
		err := emitType(f, doc, name, schema)
		if err != nil {
			return err
		}
	}
	return nil
}

func emitType(f *gen.GoGenFile, doc *v3.Document, name string, schema v3.Schema) error {
	switch schema.Type {
	case v3.String:
		fallthrough
	case v3.Number:
		fallthrough
	case v3.Integer:
		fallthrough
	case v3.Boolean:
		fallthrough
	case v3.Array:
		panic(schema)
	case v3.Object:
		return emitStruct(f, doc, name, schema)
	default:
		panic(schema.Type)
	}
	return nil
}

func emitStruct(f *gen.GoGenFile, doc *v3.Document, name string, schema v3.Schema) error {
	f.Printf(gen.Comment(schema.Description))
	f.Printf("type %s struct{\n", name)
	f.ShiftRight()
	for _, fieldName := range gen.SortedKeys(schema.Properties) {
		field := schema.Properties[fieldName]
		f.Printf(gen.Comment(field.Description))
		f.Printf("%s %s\n", gen.Public(fieldName), typeName(doc, field))
	}
	f.ShiftLeft()
	f.Printf("}\n\n")
	return nil
}

func typeName(doc *v3.Document, schema v3.Schema) string {
	switch schema.Type {
	case v3.String:
		return "string"
	case v3.Number:
		return "float64"
	case v3.Integer:
		return "int"
	case v3.Boolean:
		return "bool"
	case v3.Array:
		return "[]" + typeName(doc, *schema.Items.Schema)
	case v3.Object:
		return *schema.Ref
	default:
		if schema.Ref != nil {
			name, schema := doc.ResolveRef(*schema.Ref)
			if schema != nil {
				return name
			}
		}
		panic(fmt.Sprintf("%+v", schema))
	}
}
