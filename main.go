// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

//go:generate gotext -srclang=en update -out=catalog.go -lang=en,es

package main

import (
	"megpoid.dev/go/go-skel/cmd"
)

func main() {
	cmd.Execute()
}
