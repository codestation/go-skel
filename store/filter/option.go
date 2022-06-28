package filter

import "megpoid.xyz/go/go-skel/model/request"

// Option for filter
type Option interface {
	Apply(p *Filter)
}

// Config for filter
type Config struct {
	Rules   []Rule
	Filters []request.Filter
}

// Apply applies config to paginator
func (c *Config) Apply(p *Filter) {
	if c.Rules != nil {
		p.SetRules(c.Rules...)
	}
	// only set keys when no rules presented
	if c.Filters != nil {
		p.SetFilters(c.Filters...)
	}
}

// WithRules configures rules for query
func WithRules(rules ...Rule) Option {
	return &Config{
		Rules: rules,
	}
}

// WithFilters configures filter for query
func WithFilters(filters ...request.Filter) Option {
	return &Config{
		Filters: filters,
	}
}
