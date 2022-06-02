package web

import (
	"mtui/app"
)

type Api struct {
	app *app.App
}

func NewApi(app *app.App) *Api {
	return &Api{
		app: app,
	}
}
