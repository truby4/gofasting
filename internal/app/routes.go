package app

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (app *Application) Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(commonHeaders,
		app.recoverPanic,
		app.sessionManager.LoadAndSave,
		app.logRequest,
	)

	r.Mount("/", app.web.Routes())
	r.Mount("/api/v1/", app.api.Routes())

	return r
}
