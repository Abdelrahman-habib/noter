package ui

import "embed"

//go:embed "static/css/*.css" "static/img" "static/js" "html"
var Files embed.FS
