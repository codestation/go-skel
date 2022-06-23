// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package request

type Filter struct {
	Limit  *int    `query:"limit"`
	After  *string `query:"after"`
	Before *string `query:"before"`
}
