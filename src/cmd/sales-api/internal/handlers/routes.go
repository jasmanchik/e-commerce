package handlers

import (
	"log"
	"net/http"

	"github.com/jasmanchik/garage-sale/internal/platform/web"
	"github.com/jmoiron/sqlx"
)

func Routes(logger *log.Logger, db *sqlx.DB) http.Handler {

	app := web.NewApp(logger)

	c := Check{DB: db}
	app.Handle(http.MethodGet, "/api/health", c.Health)

	p := Product{DB: db, Log: logger}
	app.Handle(http.MethodGet, "/api/products", p.List)
	app.Handle(http.MethodPost, "/api/products", p.Create)
	app.Handle(http.MethodGet, "/api/products/{id}", p.Retrieve)
	app.Handle(http.MethodDelete, "/api/products/{id}", p.Delete)
	app.Handle(http.MethodPut, "/api/products/{id}", p.Update)

	app.Handle(http.MethodPost, "/api/products/{id}/sales", p.AddSale)
	app.Handle(http.MethodGet, "/api/products/{id}/sales", p.ListSale)

	return app
}
