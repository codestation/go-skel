// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package testdata

import "embed"

//go:embed sql
var sqlAssets embed.FS

func SqlAssets() embed.FS {
	return sqlAssets
}
