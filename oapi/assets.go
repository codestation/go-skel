// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package oapi

import "embed"

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.13.0 -package oapi -generate types -o types.gen.go goapp.yaml
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.13.0 -package oapi -generate server -o server.gen.go goapp.yaml

//go:embed goapp.yaml
var content embed.FS

func Assets() embed.FS {
	return content
}
