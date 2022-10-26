package public

import (
	"embed"
)

//go:embed js/* index.html
//go:embed node_modules/bootstrap/dist/css/bootstrap.min.css
//go:embed node_modules/vue/dist/vue.global.prod.js
//go:embed node_modules/vue-router/dist/vue-router.global.prod.js
//go:embed node_modules/vue-i18n/dist/vue-i18n.global.prod.js
//go:embed node_modules/@fortawesome/fontawesome-free/css/all.min.css
//go:embed node_modules/@fortawesome/fontawesome-free/webfonts/*
//go:embed node_modules/@vuepic/vue-datepicker/dist/main.css
//go:embed node_modules/@vuepic/vue-datepicker/dist/vue-datepicker.iife.js
var Webapp embed.FS
