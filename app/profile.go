// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package app

import (
	"context"
	"errors"
	"golang.org/x/text/message"
	"megpoid.xyz/go/go-skel/app/i18n"
	"megpoid.xyz/go/go-skel/model"
	"megpoid.xyz/go/go-skel/model/request"
	"megpoid.xyz/go/go-skel/model/response"
	"megpoid.xyz/go/go-skel/store"
)

func (a *App) GetProfile(ctx context.Context, id model.ID) (*model.Profile, error) {
	t := message.NewPrinter(i18n.GetLanguageTagsContext(ctx))

	profile, err := a.Srv().Store.Profile().Get(ctx, id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, NewAppError(t.Sprintf("Profile not found"), err)
		} else {
			return nil, NewAppError(t.Sprintf("Failed to get profile"), err)
		}
	}

	return profile, nil
}

func (a *App) ListProfiles(ctx context.Context, query *request.QueryParams) (*response.ListResponse[model.Profile], error) {
	t := message.NewPrinter(i18n.GetLanguageTagsContext(ctx))

	result, err := a.Srv().Store.Profile().List(ctx, store.WithFilter(query))
	if err != nil {
		return nil, NewAppError(t.Sprintf("Failed to list profiles"), err)
	}

	return result, nil
}
