// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package paginator

// Order type for order
type Order string

// Orders
const (
	ASC  Order = "ASC"
	DESC Order = "DESC"
)

func (o *Order) flip() Order {
	if *o == ASC {
		return DESC
	}
	return ASC
}

func (o *Order) validate() error {
	if *o != ASC && *o != DESC {
		return ErrInvalidOrder
	}
	return nil
}
