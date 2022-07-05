package sqlstore

import (
	"context"
	"megpoid.xyz/go/go-skel/model"
)

func attachRelation[T any, PT model.Modelable[T], U any, PU model.Modelable[U]](
	ctx context.Context,
	entries []PT,
	getRelationId func(m PT) model.ID,
	setRelation func(m PT, r PU),
	listByIds func(ctx context.Context, ids []model.ID) ([]PU, error),
) error {
	// exit early to don't call listByIds with an empty array
	if len(entries) == 0 {
		return nil
	}
	// list to hold the identifiers to query
	var idList []model.ID
	// map used to keep the above list with unique items
	var uniqueMap = map[model.ID]struct{}{}

	for _, entry := range entries {
		id := getRelationId(entry)
		if _, ok := uniqueMap[id]; !ok {
			uniqueMap[id] = struct{}{}
			idList = append(idList, id)
		}
	}

	results, err := listByIds(ctx, idList)
	if err != nil {
		return err
	}

	// keep the results in a map for quicker access
	var resultMap = map[model.ID]PU{}
	for _, result := range results {
		resultMap[result.GetID()] = result
	}

	for _, entry := range entries {
		if result, ok := resultMap[getRelationId(entry)]; ok {
			setRelation(entry, result)
		}
	}

	return nil
}
