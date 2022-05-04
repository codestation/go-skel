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
	"log"
	"path"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

type TranslateFunc func(id string, params map[string]interface{}) string

// T is the translation function using the default locale
var T TranslateFunc
var appBundle *i18n.Bundle

func Init() {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	assets := Assets()
	files, err := assets.ReadDir("translations")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if _, err := bundle.LoadMessageFileFS(Assets(), path.Join("translations", file.Name())); err != nil {
			panic(err)
		}
	}
	appBundle = bundle
	T = GetTranslations()
}

func GetTranslations(langs ...string) TranslateFunc {
	return func(id string, params map[string]interface{}) string {
		return TranslateLocale(id, params, langs...)
	}
}

func TranslateLocale(id string, params map[string]interface{}, langs ...string) string {
	localizer := i18n.NewLocalizer(appBundle, langs...)
	result, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: id,
		},
		TemplateData: params,
	})

	if err != nil {
		log.Printf("Cannot find translation for %s: %s", id, err)
		return id
	}

	return result
}
