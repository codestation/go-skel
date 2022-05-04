package i18n

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const LanguageTags = "LanguageTags"

func LoadMessagePrinter(preferLangKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, _ := c.Get(preferLangKey).(string)
			pref := c.Request().Header.Get("Accept-Language")
			lang := c.QueryParam("lang")
			c.Set(LanguageTags, message.MatchLanguage(lang, user, pref))
			return next(c)
		}
	}
}

func GetLanguageTags(c echo.Context) language.Tag {
	if tags, ok := c.Get(LanguageTags).(language.Tag); ok {
		return tags
	}
	return language.English
}
