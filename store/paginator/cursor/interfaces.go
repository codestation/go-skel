// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package cursor

// CustomType is an interface that custom types need to implement
// in order to allow pagination over fields inside custom types.
type CustomType interface {
	// GetCustomTypeValue returns the value corresponding to the meta attribute inside the custom type.
	GetCustomTypeValue(meta any) (any, error)
}
