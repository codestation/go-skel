// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package oapi

import "embed"

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1 -package oapi -generate types -o types.gen.go openapi.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1 -package oapi -generate server -o server.gen.go -templates templates/ openapi.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1 -package oapi -generate spec -o spec.gen.go openapi.yaml

//go:embed openapi.yaml
var content embed.FS

func Assets() embed.FS {
	return content
}
