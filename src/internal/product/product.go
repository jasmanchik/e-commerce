package product

import (
	"github.com/jmoiron/sqlx"
)

func List(db *sqlx.DB) ([]Product, error) {
	list := make([]Product, 0)
	q := "SELECT product_id, name, cost, quantity, date_created, date_updated FROM products"
	if err := db.Select(&list, q); err != nil {
		return nil, err
	}
	return list, nil
}

func Retrieve(db *sqlx.DB, id string) (*Product, error) {
	p := Product{}
	q := "SELECT product_id, name, cost, quantity, date_created, date_updated FROM products WHERE product_id = $1"
	if err := db.Get(&p, q, id); err != nil {
		return nil, err
	}
	return &p, nil
}
