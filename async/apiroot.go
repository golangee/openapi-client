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
	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}

`
