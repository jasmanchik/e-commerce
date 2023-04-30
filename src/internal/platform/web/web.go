package web

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

type App struct {
	mux *chi.Mux
	log *log.Logger
}

func NewApp(logger *log.Logger) *App {
	return &App{
		mux: chi.NewRouter(),
		log: logger,
	}
}
func (App *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	App.mux.ServeHTTP(w, r)
}

func (App *App) Handle(method, pattern string, handler http.HandlerFunc) {
	App.mux.MethodFunc(method, pattern, handler)
}
