# openapi-client
an openapi client generator fitting best for wasm integration (async callbacks)

## alternatives
There is the [official generator](https://github.com/OpenAPITools/openapi-generator), which has still a lot of
issues generating compilable and correct go code:
* The generated client looks partially unidiomatic (e.g. triple returns) 
* Cannot be safely used in wasm: typesafe and cancelable callbacks needs still to be written by hand
* No UUID support
* generates to many files for something which should never be modified by hand. 
* ugly to integrate into a versioned *go generate*