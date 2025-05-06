// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

//go:generate go run golang.org/x/text/cmd/gotext@v0.25.0 -srclang=en update -out=catalog.go -lang=en,es
//go:generate go run github.com/vektra/mockery/v3@v3.2.5
package main

import (
	"go.megpoid.dev/go-skel/cmd"
)

func main() {
	cmd.Execute()
}
