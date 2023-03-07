// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package filter

// Option for filter
type Option interface {
	Apply(f *Filter)
}

// Config for filter
type Config struct {
	Rules      []Rule
	Conditions []Condition
}

// Apply applies config to paginator
func (c *Config) Apply(f *Filter) {
	if c.Rules != nil {
		f.SetRules(c.Rules...)
	}
	// only set keys when no rules presented
	if c.Conditions != nil {
		f.SetConditions(c.Conditions...)
	}
}

// WithRules configures rules for query
func WithRules(rules ...Rule) Option {
	return &Config{
		Rules: rules,
	}
}

// WithConditions configures filter for query
func WithConditions(filters ...Condition) Option {
	return &Config{
		Conditions: filters,
	}
}
