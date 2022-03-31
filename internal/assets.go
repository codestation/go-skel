package internal

import (
	"embed"
)

//go:embed migrations
var Content embed.FS
