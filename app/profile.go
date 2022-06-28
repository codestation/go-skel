package app

import (
	"context"
	"errors"
	"golang.org/x/text/message"
	"megpoid.xyz/go/go-skel/app/i18n"
	"megpoid.xyz/go/go-skel/model"
	"megpoid.xyz/go/go-skel/model/request"
	"megpoid.xyz/go/go-skel/store"
	"megpoid.xyz/go/go-skel/store/paginator/cursor"
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

func (a *App) ListProfiles(ctx context.Context, query *request.QueryParams) ([]*model.Profile, *cursor.Cursor, error) {
	t := message.NewPrinter(i18n.GetLanguageTagsContext(ctx))

	profile, cur, err := a.Srv().Store.Profile().List(ctx, store.WithFilter(query))
	if err != nil {
		return nil, nil, NewAppError(t.Sprintf("Failed to list profiles"), err)
	}

	return profile, cur, nil
}
