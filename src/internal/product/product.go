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
