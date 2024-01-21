package public

import (
	"embed"
)

//go:embed js/* pics/* *.html
//go:embed node_modules/bootstrap/dist/css/bootstrap.min.css
//go:embed node_modules/bootswatch/dist/cyborg/bootstrap.min.css
//go:embed node_modules/bootswatch/dist/darkly/bootstrap.min.css
//go:embed node_modules/bootswatch/dist/flatly/bootstrap.min.css
//go:embed node_modules/bootswatch/dist/united/bootstrap.min.css
//go:embed node_modules/vue/dist/vue.global.js
//go:embed node_modules/vue/dist/vue.global.prod.js
//go:embed node_modules/vue-router/dist/vue-router.global.js
//go:embed node_modules/vue-router/dist/vue-router.global.prod.js
//go:embed node_modules/@fortawesome/fontawesome-free/css/all.min.css
//go:embed node_modules/@fortawesome/fontawesome-free/webfonts/*
//go:embed node_modules/@vuepic/vue-datepicker/dist/main.css
//go:embed node_modules/@vuepic/vue-datepicker/dist/vue-datepicker.iife.js
//go:embed node_modules/chart.js/dist/chart.umd.js
//go:embed node_modules/codemirror/lib/codemirror.*
//go:embed node_modules/codemirror/mode/lua/lua.js
//go:embed node_modules/codemirror/mode/javascript/javascript.js
//go:embed node_modules/codemirror/mode/htmlmixed/htmlmixed.js
//go:embed node_modules/codemirror/mode/xml/xml.js
//go:embed node_modules/codemirror/mode/css/css.js
//go:embed node_modules/codemirror/mode/toml/toml.js
var Webapp embed.FS

const DefaultCss = "node_modules/bootstrap/dist/css/bootstrap.min.css"

var ThemeMap = map[string]string{
	"default": DefaultCss,
	"cyborg":  "node_modules/bootswatch/dist/cyborg/bootstrap.min.css",
	"darkly":  "node_modules/bootswatch/dist/darkly/bootstrap.min.css",
	"flatly":  "node_modules/bootswatch/dist/flatly/bootstrap.min.css",
	"united":  "node_modules/bootswatch/dist/united/bootstrap.min.css",
}
