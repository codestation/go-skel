// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

//go:generate go run golang.org/x/text/cmd/gotext@v0.8.0 -srclang=en update -out=catalog.go -lang=en,es

package main

import (
	"go.megpoid.dev/go-skel/cmd"
)

func main() {
	cmd.Execute()
}
