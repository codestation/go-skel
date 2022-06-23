// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package paginator

import "errors"

var (
	ErrInvalidCursor = errors.New("invalid cursor for paginating")
	ErrInvalidLimit  = errors.New("limit should be greater than 0")
	ErrInvalidModel  = errors.New("model fields should match rules or keys specified for paginator")
	ErrInvalidOrder  = errors.New("order should be ASC or DESC")
	ErrNoRule        = errors.New("paginator should have at least one rule")
)
