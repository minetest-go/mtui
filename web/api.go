package web

import "mtadmin/app"

type Api struct {
	app *app.App
}

func NewApi(app *app.App) *Api {
	return &Api{app: app}
}
