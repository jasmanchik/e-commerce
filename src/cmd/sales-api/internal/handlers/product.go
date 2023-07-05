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

func (p *Product) List(w http.ResponseWriter, r *http.Request) error {
	list, err := product.List(r.Context(), p.DB)
	if err != nil {
		return err
	}
	return web.Response(w, list, http.StatusOK)
}

func (p *Product) Retrieve(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	list, err := product.Retrieve(r.Context(), p.DB, id)
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
	prod, err := product.Create(r.Context(), p.DB, &np, time.Now())
	if err != nil {
		return err
	}
	return web.Response(w, prod, http.StatusCreated)
}

func (p *Product) Delete(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	if err := product.Delete(r.Context(), p.DB, id); err != nil {
		switch err {
		case product.ErrInvalidID:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return errors.Wrapf(err, "deleting product %q", id)
		}
	}

	return web.Response(w, nil, http.StatusNoContent)
}

func (p *Product) AddSale(w http.ResponseWriter, r *http.Request) error {
	var ns product.NewSale

	if err := web.Decode(r, &ns); err != nil {
		return errors.Wrapf(err, "decoding new sale")
	}

	productID := chi.URLParam(r, "id")

	sale, err := product.AddSale(r.Context(), p.DB, ns, productID, time.Now())
	if err != nil {
		return errors.Wrapf(err, "adding new sale")
	}

	return web.Response(w, sale, http.StatusCreated)
}

func (p *Product) ListSale(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	list, err := product.ListSales(r.Context(), p.DB, id)
	if err != nil {
		return errors.Wrapf(err, "getting sale list")
	}

	return web.Response(w, list, http.StatusOK)
}

func (p *Product) Update(w http.ResponseWriter, r *http.Request) error {
	pid := chi.URLParam(r, "id")

	var update product.UpdateProduct
	if err := web.Decode(r, &update); err != nil {
		return errors.Wrap(err, "decoding product update")
	}

	if err := product.Update(r.Context(), p.DB, pid, update, time.Now()); err != nil {
		switch err {
		case product.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case product.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "updating product %q", pid)
		}
	}

	return web.Response(w, nil, http.StatusNoContent)
}
