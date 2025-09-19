package main

import (
	"net/http"

	fileserver "github.com/Abdelrahman-habib/noter/internal/file-server"
	"github.com/Abdelrahman-habib/noter/ui"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	app.logger.Debug("registering routes")
	mux := http.NewServeMux()
	fileServer := fileserver.NewFileServer(ui.Files)

	mux.Handle("GET /static/", fileServer)
	mux.HandleFunc("GET /ping", app.ping)
	dynamic := alice.New(app.sessionManager.LoadAndSave, app.noSurfMiddleware, app.authenticateMiddleware)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /about", dynamic.ThenFunc(app.about))
	mux.Handle("GET /notes", dynamic.ThenFunc(app.listNotes))
	mux.Handle("GET /note/view/{id}", dynamic.ThenFunc(app.noteView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	portected := dynamic.Append(app.requireAuthenticationMiddleware)

	mux.Handle("GET /my-notes", portected.ThenFunc(app.myNotes))
	mux.Handle("GET /note/create", portected.ThenFunc(app.noteCreate))
	mux.Handle("POST /note/create", portected.ThenFunc(app.noteCreatePost))
	mux.Handle("GET /note/edit/{id}", portected.ThenFunc(app.noteEdit))
	mux.Handle("POST /note/edit/{id}", portected.ThenFunc(app.noteCreatePost))
	mux.Handle("POST /note/delete/{id}", portected.ThenFunc(app.noteDeletePost))
	mux.Handle("GET /account/view", portected.ThenFunc(app.accountView))
	mux.Handle("GET /account/password/update", portected.ThenFunc(app.accountPasswordUpdate))
	mux.Handle("POST /account/password/update", portected.ThenFunc(app.accountPasswordUpdatePost))
	mux.Handle("POST /user/logout", portected.ThenFunc(app.userLogoutPost))

	app.logger.Debug("routes registered")

	standard := alice.New(app.recoverPanicMiddleware, app.loggerMiddleware, commonHeadersMiddleware)

	return standard.Then(mux)
}
