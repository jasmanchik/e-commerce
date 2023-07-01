package handlers

import (
	"log"
	"net/http"

	"github.com/jasmanchik/garage-sale/internal/platform/web"
	"github.com/jmoiron/sqlx"
)

func Routes(logger *log.Logger, db *sqlx.DB) http.Handler {

	app := web.NewApp(logger)
	p := Product{DB: db, Log: logger}

	app.Handle(http.MethodGet, "/api/products", p.List)
	app.Handle(http.MethodPost, "/api/products", p.Create)
	app.Handle(http.MethodGet, "/api/products/{id}", p.Retrieve)
	app.Handle(http.MethodDelete, "/api/products/{id}", p.Delete)

	return app
}
