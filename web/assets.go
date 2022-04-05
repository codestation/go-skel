package web

import "embed"

//go:embed static
var assets embed.FS

func Assets() embed.FS {
	return assets
}
