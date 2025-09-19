package main

import (
	"crypto/tls"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"

	"github.com/Abdelrahman-habib/noter/internal/models"
)

type application struct {
	logger         *slog.Logger
	config         *config
	notes          models.NoteModelInterface
	users          models.UserModelInterface
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func (app *application) serve() error {
	// only elliptic curves with assembly implementations are used
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	server := &http.Server{
		Addr:    app.config.addr,
		Handler: app.routes(),

		ErrorLog:  slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
		TLSConfig: tlsConfig,

		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Minute,
	}
	app.logger.Info("starting server", slog.String("addr", server.Addr))

	if app.config.tlsCert != "" && app.config.tlsKey != "" {
		return server.ListenAndServeTLS(app.config.tlsCert, app.config.tlsKey)
	}
	return server.ListenAndServe()

}
