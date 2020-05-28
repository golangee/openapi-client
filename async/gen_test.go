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

import "testing"

func TestGenerate(t *testing.T) {
	err := Generate([]byte(spec), Options{
		TargetDir:     "",
		TargetPackage: "blub",
	})

	if err != nil {
		t.Fatal(err)
	}
}

const spec = `{
   "openapi":"3.0.1",
   "info":{
      "title":"test",
      "version":"",
	  "description":"test is for testing"
   },
   "servers":[
      {
         "url":"http://localhost:8080"
      }
   ],
   "paths":{
      "/api/v1/setup/status":{
         "get":{
            "tags":[
               "setup"
            ],
            "summary":"Status returns the current setup status",
            "description":"Status returns the current setup status. This is usually only relevant in the installation phase.",
            "responses":{
               "200":{
                  "description":"Status represents the current setup status.",
                  "content":{
                     "application/json":{
                        "schema":{
                           "type":"array",
                           "items":{
                              "$ref":"#/components/schemas/Status"
                           }
                        }
                     }
                  }
               }
            }
         }
      }
   },
   "components":{
      "schemas":{
         "Status":{
            "type":"object",
            "properties":{
               "Id":{
                  "type":"integer",
                  "format":"int32",
                  "description":"status id\n"
               },
               "Message":{
                  "type":"string",
                  "description":"a textual representation as a developer notice\n"
               }
            },
            "description":"Status represents the current setup status.\n",
            "x-ee.type":"github.com/worldiety/test/internal/service/setup#Status"
         }
      }
   }
}
`
