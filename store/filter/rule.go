// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package filter

type VariableType string

const (
	VariableString          VariableType = "string"
	VariableInteger         VariableType = "number"
	VariableBool            VariableType = "boolean"
	VariableDecimal         VariableType = "decimal"
	VariableDate            VariableType = "date"
	VariableTimestamp       VariableType = "timestamp"
	VariableTimestampMillis VariableType = "timestamp_millis"
)

type OperationType string

const (
	OperationEqual          OperationType = "eq"
	OperationNotEqual       OperationType = "neq"
	OperationLessThan       OperationType = "lt"
	OperationLessOrEqual    OperationType = "lte"
	OperationGreaterThan    OperationType = "gt"
	OperationGreaterOrEqual OperationType = "gte"
	OperationHas            OperationType = "has"
	OperationIn             OperationType = "in"
	OperationIsNull         OperationType = "isnull"
)

type Rule struct {
	Key        string
	Operation  []OperationType
	Type       VariableType
	AcceptNull bool
}
