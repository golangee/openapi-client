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
	"github.com/golangee/openapi-client/internal/gen"
	v3 "github.com/golangee/openapi/v3"
	"regexp"
	"strings"
)

type endpoint struct {
	path   string
	method string
	op     *v3.Operation
}

func emitCallGroups(f *gen.GoGenFile, parentType string, doc *v3.Document) error {
	groups := map[string][]endpoint{}
	for path, call := range doc.Paths {
		for method, op := range call.Map() {
			tmpTags := []string{doc.Info.Title}
			if len(op.Tags) > 0 {
				tmpTags = op.Tags
			}

			for _, tag := range tmpTags {
				list := groups[tag]
				list = append(list, endpoint{path, method, op})
				groups[tag] = list
			}

		}
	}

	for _, tag := range gen.SortedKeys(groups) {
		endpoints := groups[tag]
		err := emitCallGroup(f, doc, parentType, gen.Public(tag)+"Service", endpoints)
		if err != nil {
			return err
		}
	}

	return nil
}

func emitCallGroup(f *gen.GoGenFile, doc *v3.Document, parentType string, name string, endpoints []endpoint) error {
	f.Printf("// %s groups tagged api calls\n", name)
	f.Printf("type %s struct {\n", name)
	f.Printf("parent *%s\n", parentType)
	f.Printf("}\n\n")
	for _, ep := range endpoints {
		err := emitSyncCall(f, doc, name, ep)
		if err != nil {
			return err
		}

		err = emitAsyncCall(f, doc, name, ep)
		if err != nil {
			return err
		}
	}
	return nil
}

func emitSyncCall(f *gen.GoGenFile, doc *v3.Document, receiverTypeName string, ep endpoint) error {
	f.Printf(gen.Comment(ep.op.Description))
	f.Printf("func (_self %s) sync%s(_ctx %s) (%s,error){\n", receiverTypeName, methodName(ep), f.ImportName("context", "Context"), pickResponseAndResolveTypeName(doc, ep))
	pathParams := pathParamsToSprintf(ep.path)
	f.Printf("url := %s(\"%s\",%s)\n", f.ImportName("fmt", "Sprintf"), pathParams.sprintfPath, strings.Join(pathParams.params, ","))
	f.Printf("}\n")
	return nil
}

func emitAsyncCall(f *gen.GoGenFile, doc *v3.Document, receiverTypeName string, ep endpoint) error {
	f.Printf(gen.Comment(ep.op.Description))
	f.Printf("func (_self %s) %s(_ctx %s, f func(res %s,err error)){\n", receiverTypeName, methodName(ep), f.ImportName("context", "Context"), pickResponseAndResolveTypeName(doc, ep))
	f.Printf("go func(){\n")
	f.Printf("res,err := %s(_ctx)\n", "sync"+methodName(ep))
	f.Printf("f(res,err)\n")
	f.Printf("}()\n")
	f.Printf("}\n")
	return nil
}

func methodName(ep endpoint) string {
	path := strings.ReplaceAll(ep.path, "{", "")
	path = strings.ReplaceAll(path, "}", "")
	method := gen.Public(strings.ToLower(ep.method))
	if method == "Get" {
		method = ""
	}

	if method == "Post" && strings.Contains(strings.ToLower(ep.op.Summary), "create") {
		method = "Create"
	}

	return gen.SlashToCamelCase(method + "/" + path)
}

func pickResponseAndResolveTypeName(doc *v3.Document, ep endpoint) string {
	response, has200 := ep.op.Responses["200"]
	if has200 {
		json, hasJson := response.Content["application/json"]
		if hasJson {
			return typeName(doc, json.Schema)
		}
	}
	return "interface{}"
}

type namedPath struct {
	path        string
	sprintfPath string
	params      []string
}

func pathParamsToSprintf(path string) namedPath {
	r := namedPath{path: path}
	regex := regexp.MustCompile(`{\w*}`)
	sprint := regex.ReplaceAllStringFunc(path, func(s string) string {
		name := s[1 : len(s)-1]
		r.params = append(r.params, name)
		return "%v"
	})
	r.sprintfPath = sprint
	return r
}
