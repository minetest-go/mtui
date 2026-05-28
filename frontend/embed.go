package frontend

import (
	"embed"
)

//go:embed dist/*
var Webapp embed.FS
