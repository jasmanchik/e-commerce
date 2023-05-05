package product

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Predefined errors for known failure scenarios
var (
	ErrNotFound  = errors.New("product not found")
	ErrInvalidID = errors.New("id provided was not a valid UUID")
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

	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}

	p := Product{}
	q := "SELECT product_id, name, cost, quantity, date_created, date_updated FROM products WHERE product_id = $1"
	if err := db.Get(&p, q, id); err != nil {
		return nil, ErrNotFound
	}
	return &p, nil
}

func Create(db *sqlx.DB, np *NewProduct, now time.Time) (*Product, error) {
	p := Product{
		ID:          uuid.New().String(),
		Name:        np.Name,
		Cost:        np.Cost,
		Quantity:    np.Quantity,
		DateCreated: now,
		DateUpdated: now,
	}
	const q = "INSERT INTO products (product_id, name, cost, quantity, date_created, date_updated) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := db.Exec(q, p.ID, p.Name, p.Cost, p.Quantity, p.DateCreated, p.DateUpdated)
	if err != nil {
		return nil, errors.Wrapf(err, "inserting product: %v", np)
	}
	return &p, nil
}
