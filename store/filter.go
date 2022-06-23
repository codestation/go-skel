// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package store

import (
	"megpoid.xyz/go/go-skel/model/request"
	"megpoid.xyz/go/go-skel/store/paginator"
)

type FilterOption func(paginator *paginator.Paginator)

func WithFilter(query *request.QueryParams) FilterOption {
	return func(paginator *paginator.Paginator) {
		if query.Pagination.Limit != nil {
			paginator.SetLimit(*query.Pagination.Limit)
		}
		if query.Pagination.After != nil {
			paginator.SetAfterCursor(*query.Pagination.After)
		}
		if query.Pagination.Before != nil {
			paginator.SetBeforeCursor(*query.Pagination.Before)
		}
	}
}
