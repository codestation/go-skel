// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package usecase

import (
	"context"

	"golang.org/x/text/message"
	"megpoid.dev/go/go-skel/pkg/i18n"
)

// common has all the repositories and other data required for all the usecases.
type common struct{}

// printer returns a printer to localize messages to other languages.
func (u *common) printer(ctx context.Context) *message.Printer {
	return message.NewPrinter(i18n.GetLanguageTagsContext(ctx))
}

func newCommon() common {
	return common{}
}
