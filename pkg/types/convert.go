// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package types

func AsPointer[T any](v T) *T {
	return &v
}

func AsValue[T any](v *T) T {
	if v != nil {
		return *v
	}
	var empty T
	return empty
}
