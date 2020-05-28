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
	"io/ioutil"
	"os"
	"path/filepath"
)

// Options to use for generating a new client
type Options struct {
	TargetDir     string
	TargetPackage string
}

// Generates determines the root of the module and applies the options to generate a new client from the spec.
func Generate(spec []byte, opts Options) error {
	doc, err := v3.FromJson(spec)
	if err != nil {
		return fmt.Errorf("unable to parse document: %w", err)
	}

	file := gen.NewGoGenFile(opts.TargetPackage, "openapi-client")

	err = emitTypes(file, doc)
	if err != nil {
		return fmt.Errorf("unable to emit types: %w", err)
	}

	parentType, err := emitApiRoot(file, doc)
	if err != nil {
		return fmt.Errorf("unable to emit api root: %w", err)
	}

	err = emitCallGroups(file, parentType, doc)
	if err != nil {
		return fmt.Errorf("unable to emit call groups: %w", err)
	}

	fmt.Println(file.FormatString())

	dir, err := gen.ModRootDir()
	if err != nil {
		return err
	}

	fname := filepath.Join(dir, opts.TargetDir, "openapiclient.gen.go")
	if err := ioutil.WriteFile(fname, []byte(file.FormatString()), os.ModePerm); err != nil {
		return err
	}

	return nil
}
