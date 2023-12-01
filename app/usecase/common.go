// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package usecase

import (
	"context"
	"time"

	"go.megpoid.dev/go-skel/pkg/i18n"
	"golang.org/x/text/message"
)

// common has all the repositories and other data required for all the usecases.
type common struct {
	timeNow func() time.Time
}

// printer returns a printer to localize messages to other languages.
func (u *common) printer(ctx context.Context) *message.Printer {
	return message.NewPrinter(i18n.GetLanguageTagsContext(ctx))
}

func (u *common) currentTime() time.Time {
	return u.timeNow()
}

type Option func(m *common)

func WithTime(timeFn func() time.Time) Option {
	return func(m *common) {
		m.timeNow = timeFn
	}
}

func newCommon(opts ...Option) common {
	c := common{}
	for _, opt := range opts {
		opt(&c)
	}

	if c.timeNow == nil {
		c.timeNow = time.Now
	}

	return c
}
