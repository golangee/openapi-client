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

package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"unicode"
)

// ModRootDir returns the root directory of current module. If the current working directory is not a module
// returns an error.
func ModRootDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	root := cwd
	for {
		stat, err := os.Stat(filepath.Join(root, "go.mod"))
		if err == nil && stat.Mode().IsRegular() {
			return root, nil
		}
		root = filepath.Dir(root)
		if root == "/" || root == "." {
			return "", fmt.Errorf("%s is not withing a go module", cwd)
		}
	}
}

// Public ensures that str starts with an uppercase letter
func Public(str string) string {
	if str == "" {
		return str
	}
	return string(unicode.ToUpper(rune(str[0]))) + str[1:]
}

// Comment assembles a string with correct newlines and // at the beginning of each line
func Comment(str string) string {
	str = strings.TrimSpace(str)
	sb := &strings.Builder{}
	for _, line := range strings.Split(str, "\n") {
		sb.WriteString("// ")
		sb.WriteString(line)
		sb.WriteString("\n")
	}
	return sb.String()
}

// SortedKeys is missing generics very badly, m is actually map[string]<T>
func SortedKeys(m interface{}) []string {
	keys := reflect.ValueOf(m).MapKeys()
	res := make([]string, len(keys))
	for i, v := range keys {
		res[i] = v.String()
	}
	sort.Strings(res)
	return res
}

// SlashToCamelCase makes from /my/path a string like MyPath
func SlashToCamelCase(str string) string {
	sb := &strings.Builder{}
	nextUp := true
	for _, r := range str {
		if r == '/' {
			nextUp = true
			continue
		}
		if nextUp {
			sb.WriteRune(unicode.ToUpper(r))
			nextUp = false
		} else {
			sb.WriteRune(r)
		}

	}
	return sb.String()
}
