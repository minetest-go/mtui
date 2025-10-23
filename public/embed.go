package public

import (
	"embed"
)

//go:embed js/* pics/* *.html
//go:embed wasm/*
//go:embed node_modules/bootstrap/dist/css/bootstrap.min.css
//go:embed node_modules/vue/dist/vue.global.js
//go:embed node_modules/vue/dist/vue.global.prod.js
//go:embed node_modules/vue-router/dist/vue-router.global.js
//go:embed node_modules/vue-router/dist/vue-router.global.prod.js
//go:embed node_modules/@fortawesome/fontawesome-free/css/all.min.css
//go:embed node_modules/@fortawesome/fontawesome-free/webfonts/*
//go:embed node_modules/@vuepic/vue-datepicker/dist/main.css
//go:embed node_modules/@vuepic/vue-datepicker/dist/vue-datepicker.iife.js
//go:embed node_modules/codemirror/lib/codemirror.*
//go:embed node_modules/codemirror/mode/lua/lua.js
//go:embed node_modules/codemirror/mode/javascript/javascript.js
//go:embed node_modules/codemirror/mode/htmlmixed/htmlmixed.js
//go:embed node_modules/codemirror/mode/xml/xml.js
//go:embed node_modules/codemirror/mode/css/css.js
//go:embed node_modules/codemirror/mode/toml/toml.js
//go:embed node_modules/marked/lib/marked.umd.js
//go:embed node_modules/dompurify/dist/purify.min.js
var Webapp embed.FS
