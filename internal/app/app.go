package app

import (
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/truby4/go-fasting/internal/api"
	"github.com/truby4/go-fasting/internal/web"
)

type Application struct {
	Logger *log.Logger

	web *web.Handler
	api *api.Handler
}

func New() Application {
	logger := log.NewWithOptions(os.Stderr, log.Options{ReportTimestamp: true})

	return Application{
		Logger: logger,
		web:    web.NewHandler(logger),
		api:    api.NewHandler(logger),
	}
}

func (app *Application) Serve(addr string) error {
	app.Logger.Infof("Starting server on localhost %s", addr)
	srv := &http.Server{
		Addr:         addr,
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     app.Logger.StandardLog(log.StandardLogOptions{ForceLevel: log.ErrorLevel}),
	}

	err := srv.ListenAndServe()
	return err
}
