package main

import (
	"net/http"

	fileserver "github.com/Abdelrahman-habib/snippetbox/internal/file-server"
	"github.com/Abdelrahman-habib/snippetbox/ui"
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
	mux.Handle("GET /snippets", dynamic.ThenFunc(app.listSnippets))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	portected := dynamic.Append(app.requireAuthenticationMiddleware)

	mux.Handle("GET /my-snippets", portected.ThenFunc(app.mySnippets))
	mux.Handle("GET /snippet/create", portected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", portected.ThenFunc(app.snippetCreatePost))
	mux.Handle("GET /snippet/edit/{id}", portected.ThenFunc(app.snippetEdit))
	mux.Handle("POST /snippet/edit/{id}", portected.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /snippet/delete/{id}", portected.ThenFunc(app.snippetDeletePost))
	mux.Handle("GET /account/view", portected.ThenFunc(app.accountView))
	mux.Handle("GET /account/password/update", portected.ThenFunc(app.accountPasswordUpdate))
	mux.Handle("POST /account/password/update", portected.ThenFunc(app.accountPasswordUpdatePost))
	mux.Handle("POST /user/logout", portected.ThenFunc(app.userLogoutPost))

	app.logger.Debug("routes registered")

	standard := alice.New(app.recoverPanicMiddleware, app.loggerMiddleware, commonHeadersMiddleware)

	return standard.Then(mux)
}
