// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package filter

import "strings"

type Condition struct {
	Field     string
	Operation OperationType
	Value     string
	SkipRule  bool
}

func (f Condition) Values() []string {
	return strings.Split(f.Value, ",")
}
