// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package filter

import "strings"

type Condition struct {
	Field     string
	Operation OperationType
	Value     string
}

func NewCondition(field, operation, value string) Condition {
	return Condition{
		Field:     field,
		Operation: OperationType(operation),
		Value:     value,
	}
}

func (f Condition) Values() []string {
	return strings.Split(f.Value, ",")
}
