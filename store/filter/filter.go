// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package filter

import (
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/shopspring/decimal"
	"megpoid.xyz/go/go-skel/model/request"
	"strconv"
	"strings"
	"time"
)

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

// New creates a query filter
func New(opts ...Option) *Filter {
	p := &Filter{}
	for _, opt := range opts {
		opt.Apply(p)
	}
	return p
}

type Filter struct {
	rules   map[string]Rule
	filters []request.Filter
}

// SetRules sets paging rules
func (f *Filter) SetRules(rules ...Rule) {
	f.rules = make(map[string]Rule, len(rules))
	for _, rule := range rules {
		f.rules[rule.Key] = rule
	}
}

// SetFilters sets filter rules
func (f *Filter) SetFilters(filters ...request.Filter) {
	f.filters = make([]request.Filter, len(filters))
	copy(f.filters, filters)
}

type Rule struct {
	Key        string
	Operation  []request.OperationType
	Type       VariableType
	AcceptNull bool
}

func (f *Filter) getValueFromFilter(rule Rule, filter request.Filter) (any, error) {
	var value any

	if rule.AcceptNull && strings.ToLower(filter.Value) == "null" {
		return nil, nil
	}

	switch rule.Type {
	case VariableString:
		value = filter.Value
	case VariableInteger:
		intVal, err := strconv.ParseInt(filter.Value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid filter value for %s, must be integer: %w", filter.Field, err)
		}
		value = intVal
	case VariableDecimal:
		decVal, err := decimal.NewFromString(filter.Value)
		if err != nil {
			return nil, fmt.Errorf("invalid filter value for %s, must be decimal: %w", filter.Field, err)
		}
		value = decVal
	case VariableDate:
		_, err := time.Parse("2006-01-02", filter.Value)
		if err != nil {
			return nil, fmt.Errorf("invalid filter value for %s, must match format yyyy-MM-dd: %w", filter.Field, err)
		}
		value = filter.Value
	case VariableTimestamp:
		i, err := strconv.ParseInt(filter.Value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid filter value for %s, must be a timestamp: %w", filter.Field, err)
		}
		value = time.Unix(i, 0)
	case VariableTimestampMillis:
		i, err := strconv.ParseInt(filter.Value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid filter value for %s, must be a timestamp with milliseconds: %w", filter.Field, err)
		}
		value = time.UnixMilli(i)
	case VariableBool:
		boolVal, err := strconv.ParseBool(filter.Value)
		if err != nil {
			return nil, fmt.Errorf("invalid filter value for %s, must be boolean: %w", filter.Field, err)
		}
		value = boolVal
	default:
		return nil, fmt.Errorf("unknown rule type: %s", rule.Type)
	}

	return value, nil
}

func (f *Filter) Apply(query *goqu.SelectDataset) (*goqu.SelectDataset, error) {
	expr, err := f.buildWhereExpression()
	if err != nil {
		return nil, fmt.Errorf("failed to build query expression: %w", err)
	}
	return query.Where(expr), nil
}

func (f *Filter) buildWhereExpression() (exp.ExpressionList, error) {
	queries := make([]exp.Expression, 0)
	for _, filter := range f.filters {
		rule, ok := f.rules[filter.Field]
		if !ok {
			continue
		}

		value, err := f.getValueFromFilter(rule, filter)
		if err != nil {
			return nil, err
		}

		if len(rule.Operation) > 0 {
			found := false
			for _, operator := range rule.Operation {
				if operator == filter.Operation {
					found = true
					break
				}
			}

			if !found {
				return nil, fmt.Errorf("operator not permitted for field %s: %s", filter.Field, filter.Operation)
			}
		}

		var queryFilter exp.Expression
		switch filter.Operation {
		case request.OperationEqual:
			queryFilter = goqu.I(rule.Key).Eq(value)
		case request.OperationNotEqual:
			queryFilter = goqu.I(rule.Key).Neq(value)
		case request.OperationGreaterThan:
			queryFilter = goqu.I(rule.Key).Gt(value)
		case request.OperationGreaterOrEqual:
			queryFilter = goqu.I(rule.Key).Gte(value)
		case request.OperationLessThan:
			queryFilter = goqu.I(rule.Key).Lt(value)
		case request.OperationLessOrEqual:
			queryFilter = goqu.I(rule.Key).Lte(value)
		case request.OperationHas:
			queryFilter = goqu.I(rule.Key).ILike(fmt.Sprintf("%%%s%%", value))
		case request.OperationIn:
			values := strings.Split(value.(string), ",")
			queryFilter = goqu.I(rule.Key).In(values)
		}
		queries = append(queries, queryFilter)
	}

	return goqu.And(queries...), nil
}
