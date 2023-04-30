package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jasmanchik/garage-sale/internal/product"
	"github.com/jmoiron/sqlx"
)

type Product struct {
	DB *sqlx.DB
}

func (p *Product) List(w http.ResponseWriter, r *http.Request) {

	list, err := product.List(p.DB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("ListProducts: error marshalling data: %s", err)
		return
	}
	data, err := json.Marshal(list)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("ListProducts: error marshalling data: %s", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		log.Printf("ListProducts: write response: %s", err)
	}
}
