// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package cursor

import "errors"

var (
	ErrInvalidCursor = errors.New("invalid cursor")
	ErrInvalidModel  = errors.New("invalid entity")
)
