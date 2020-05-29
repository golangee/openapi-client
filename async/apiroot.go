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
)

func emitApiRoot(f *gen.GoGenFile, doc *v3.Document) (string, error) {
	f.ImportName("context", "")
	f.ImportName("encoding/json", "")
	f.ImportName("io", "")
	f.ImportName("net/http", "")
	f.ImportName("net/url", "")
	f.ImportName("io/ioutil", "")
	f.ImportName("reflect", "")
	//f.ImportName("log","")

	f.Printf(errTypeStub)
	f.Printf("\n")
	f.Printf(gen.Comment(doc.Info.Description))
	rootName := gen.Public(doc.Info.Title + "Service")
	f.Printf(parentClientStub, rootName)
	return rootName, nil
}

const parentClientStub = `// %[1]s is a basic http client implementation, which provides some reasonable defaults
type %[1]s struct {
	baseURL    *url.URL
	userAgent  string
	httpClient *http.Client
}

// New%[1]s creates a new service instance. If httpClient is nil, the default client is used.
func New%[1]s(baseURL *url.URL, userAgent string, httpClient *http.Client) *%[1]s {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &%[1]s{baseURL: baseURL, httpClient: httpClient, userAgent: userAgent}
}

func (s *%[1]s) newRequest(ctx context.Context, method, path, contentType, accept string, body io.Reader) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := s.baseURL.ResolveReference(rel)
	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", contentType)
	}

	req.Header.Set("Accept", accept)
	if s.userAgent != "" {
		req.Header.Set("User-Agent", s.userAgent)
	}
	return req, nil
}

func (s *%[1]s) doJson(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		err = json.NewDecoder(resp.Body).Decode(v)
		return resp, err
	}else{
		err := ParseError(resp.Body)
		return resp, err
	}
	
}


`

// copied from http/error.go
const errTypeStub = "// Error describes a (nested) server error\ntype Error struct {\n\tId               string      `json:\"id\"`                         // Id is unique for a specific error, e.g. mydomain.not.assigned\n\tMessage          string      `json:\"message\"`                    // Message is a string for the developer\n\tLocalizedMessage string      `json:\"localizedMessage,omitempty\"` // LocalizedMessage is something to display the user\n\tCausedBy         *Error      `json:\"causedBy,omitempty\"`         // CausedBy returns an optional root error\n\tType             string      `json:\"type,omitempty\"`             // Type is a developer notice for the internal inspection\n\tDetails          interface{} `json:\"details,omitempty\"`          // Details contains arbitrary payload\n}\n\n// ParseError tries to parse the response as json. In any case it returns an error.\nfunc ParseError(reader io.Reader) *Error {\n\tbuf, err := ioutil.ReadAll(reader)\n\tif err != nil {\n\t\treturn AsError(err)\n\t}\n\n\tres := &Error{}\n\terr = json.Unmarshal(buf, res)\n\tif err != nil {\n\t\treturn AsError(err)\n\t}\n\n\treturn res\n}\n\n// ID returns the unique error class id\nfunc (c *Error) ID() string {\n\treturn c.Id\n}\n\n// Error returns the message\nfunc (c *Error) Error() string {\n\treturn c.Message\n}\n\n// LocalizedError is like Error but translated or empty\nfunc (c *Error) LocalizedError() string {\n\treturn c.LocalizedMessage\n}\n\n// Class returns the technical type\nfunc (c *Error) Class() string {\n\treturn c.Type\n}\n\n// Payload returns the details\nfunc (c *Error) Payload() interface{} {\n\treturn c.Details\n}\n\n// Unwrap returns the cause or nil\nfunc (c *Error) Unwrap() error {\n\tif c.CausedBy == nil { // otherwise error iface will not be nil, because of the type info in interface\n\t\treturn nil\n\t}\n\treturn c.CausedBy\n}\n\nfunc AsError(err error) *Error {\n\tif e, ok := err.(*Error); ok {\n\t\treturn e\n\t}\n\n\te := &Error{}\n\te.Type = reflect.TypeOf(err).String()\n\te.Message = err.Error()\n\n\tif code, ok := err.(interface{ ID() string }); ok {\n\t\te.Id = code.ID()\n\t} else {\n\t\te.Id = e.Type\n\t}\n\n\tif details, ok := err.(interface{ Payload() interface{} }); ok {\n\t\te.Details = details\n\t}\n\n\tif localized, ok := err.(interface{ LocalizedError() string }); ok {\n\t\te.LocalizedMessage = localized.LocalizedError()\n\t}\n\n\tif class, ok := err.(interface{ Class() string }); ok {\n\t\te.Type = class.Class()\n\t}\n\n\tif wrapper, ok := err.(interface{ Unwrap() error }); ok {\n\t\tcause := wrapper.Unwrap()\n\t\tif cause != nil {\n\t\t\ttmp := AsError(cause)\n\t\t\te.CausedBy = tmp\n\t\t}\n\t}\n\n\treturn e\n}"