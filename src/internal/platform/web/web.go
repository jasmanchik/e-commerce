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

type Handler func(w http.ResponseWriter, r *http.Request) error

func NewApp(logger *log.Logger) *App {
	return &App{
		mux: chi.NewRouter(),
		log: logger,
	}
}
func (App *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	App.mux.ServeHTTP(w, r)
}

func (App *App) Handle(method, pattern string, handler Handler) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			App.log.Printf("ERROR: %v", err)

			if err := RespondError(w, err); err != nil {
				App.log.Printf("ERROR: %v", err)
			}
		}
	}
	App.mux.MethodFunc(method, pattern, fn)
}
