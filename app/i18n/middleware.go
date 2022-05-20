// Copyright (c) 2022 codestation
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package i18n

import (
	"context"

	"github.com/labstack/echo/v4"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type LanguageTagKey struct{}

func (l LanguageTagKey) String() string {
	return "LanguageTag"
}

func LoadMessagePrinter(preferLangKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, _ := c.Get(preferLangKey).(string)
			pref := c.Request().Header.Get("Accept-Language")
			lang := c.QueryParam("lang")
			tags := message.MatchLanguage(lang, user, pref)
			c.Set(LanguageTagKey{}.String(), tags)

			schemaCtx := context.WithValue(c.Request().Context(), LanguageTagKey{}, tags)
			request := c.Request().WithContext(schemaCtx)
			c.SetRequest(request)

			return next(c)
		}
	}
}

func GetLanguageTags(c echo.Context) language.Tag {
	if tags, ok := c.Get(LanguageTagKey{}.String()).(language.Tag); ok {
		return tags
	}
	return language.English
}

func GetLanguageTagsContext(ctx context.Context) language.Tag {
	if tags, ok := ctx.Value(LanguageTagKey{}).(language.Tag); ok {
		return tags
	}
	return language.English
}
