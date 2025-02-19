// Package oapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package oapi

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"LwAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
