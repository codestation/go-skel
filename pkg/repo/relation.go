// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package repo

import (
	"context"

	"megpoid.dev/go/go-skel/pkg/model"
	"megpoid.dev/go/go-skel/pkg/response"
)

func AttachRelation[T, U model.Modelable](
	ctx context.Context,
	entries []T,
	getRelationId func(m T) *int64,
	setRelation func(m T, r U),
	listByIds func(ctx context.Context, ids []int64) (*response.ListResponse[U], error),
) error {
	if len(entries) == 0 {
		return nil
	}
	// list to hold the identifiers to query
	var idList []int64
	// map used to keep the above list with unique items
	var uniqueMap = map[int64]struct{}{}

	for _, entry := range entries {
		id := getRelationId(entry)
		if id != nil {
			if _, ok := uniqueMap[*id]; !ok {
				uniqueMap[*id] = struct{}{}
				idList = append(idList, *id)
			}
		}
	}

	if idList == nil {
		return nil
	}

	results, err := listByIds(ctx, idList)
	if err != nil {
		return err
	}

	// keep the results in a map for quicker access
	var resultMap = map[int64]U{}
	for _, result := range results.Items {
		resultMap[result.GetID()] = result
	}

	for _, entry := range entries {
		id := getRelationId(entry)
		if id != nil {
			if result, ok := resultMap[*id]; ok {
				setRelation(entry, result)
			}
		}
	}

	return nil
}
