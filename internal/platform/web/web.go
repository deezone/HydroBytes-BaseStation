package web

import (
	// Core packages
	"log"
	"net/http"

	// Third-party packages
	"github.com/go-chi/chi"
)

// App is the entrypoint into the application and what controls the context of each request.
type App struct {
	log *log.Logger
	mux *chi.Mux
}

// NewApp constructs an App to handle a set of routes.
func NewApp(log *log.Logger) *App {
	return &App{
		log: log,
		mux: chi.NewRouter(),
	}
}

// Handle associates a handler function with an HTTP Method and URL pattern.
func (a *App) Handle(method, url string, h http.HandlerFunc) {
	a.mux.MethodFunc(method, url, h)
}

// ServeHTTP implements the http.Handler interface.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
