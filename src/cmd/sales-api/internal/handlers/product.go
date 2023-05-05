package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/jasmanchik/garage-sale/internal/platform/web"
	"github.com/jasmanchik/garage-sale/internal/product"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Product struct {
	DB  *sqlx.DB
	Log *log.Logger
}

func (p *Product) List(w http.ResponseWriter, _ *http.Request) error {

	list, err := product.List(p.DB)
	if err != nil {
		return err
	}
	return web.Response(w, list, http.StatusOK)
}

func (p *Product) Retrieve(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	list, err := product.Retrieve(p.DB, id)
	if err != nil {
		switch err {
		case product.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case product.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return errors.Wrapf(err, "looking for product %q", id)
		}
	}
	return web.Response(w, list, http.StatusOK)
}

func (p *Product) Create(w http.ResponseWriter, r *http.Request) error {
	var np product.NewProduct
	if err := web.Decode(r, &np); err != nil {
		return err
	}
	prod, err := product.Create(p.DB, &np, time.Now())
	if err != nil {
		return err
	}
	return web.Response(w, prod, http.StatusCreated)
}
