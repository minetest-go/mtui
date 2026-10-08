package public

import (
	"embed"
)

// built webapp (npm run build), "all:" to include files starting with "_" or "."
//
//go:embed all:dist
var Webapp embed.FS

// path prefix of the webapp files in the embedded filesystem
const Prefix = "dist"
