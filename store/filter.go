// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package store

import (
	"megpoid.xyz/go/go-skel/model/request"
	"megpoid.xyz/go/go-skel/store/filter"
	"megpoid.xyz/go/go-skel/store/paginator"
)

type FilterOption func(paginator *paginator.Paginator, filter *filter.Filter)

func WithFilter(query *request.QueryParams) FilterOption {
	return func(p *paginator.Paginator, f *filter.Filter) {
		if query.Pagination.Limit != nil {
			p.SetLimit(*query.Pagination.Limit)
		}
		if query.Pagination.After != nil {
			p.SetAfterCursor(*query.Pagination.After)
		}
		if query.Pagination.Before != nil {
			p.SetBeforeCursor(*query.Pagination.Before)
		}
		if query.Filters != nil {
			var conditions []filter.Condition
			for _, fi := range query.Filters {
				conditions = append(conditions, filter.NewCondition(fi.Field, fi.Operation, fi.Value))
			}
			f.SetConditions(conditions...)
		}
	}
}
