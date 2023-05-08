package product

import (
	"context"
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

func List(ctx context.Context, db *sqlx.DB) ([]Product, error) {
	list := make([]Product, 0)
	q := "SELECT product_id, name, cost, quantity, date_created, date_updated FROM products"
	if err := db.SelectContext(ctx, &list, q); err != nil {
		return nil, err
	}
	return list, nil
}

func Retrieve(ctx context.Context, db *sqlx.DB, id string) (*Product, error) {

	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}

	p := Product{}
	q := "SELECT product_id, name, cost, quantity, date_created, date_updated FROM products WHERE product_id = $1"
	if err := db.GetContext(ctx, &p, q, id); err != nil {
		return nil, ErrNotFound
	}
	return &p, nil
}

func Create(ctx context.Context, db *sqlx.DB, np *NewProduct, now time.Time) (*Product, error) {
	p := Product{
		ID:          uuid.New().String(),
		Name:        np.Name,
		Cost:        np.Cost,
		Quantity:    np.Quantity,
		DateCreated: now.UTC(),
		DateUpdated: now.UTC(),
	}
	const q = "INSERT INTO products (product_id, name, cost, quantity, date_created, date_updated) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := db.ExecContext(ctx, q, p.ID, p.Name, p.Cost, p.Quantity, p.DateCreated, p.DateUpdated)
	if err != nil {
		return nil, errors.Wrapf(err, "inserting product: %v", np)
	}
	return &p, nil
}

func Delete(ctx context.Context, db *sqlx.DB, p *Product) error {
	const q = "DELETE FROM products WHERE product_id=$1"
	_, err := db.ExecContext(ctx, q, p.ID)
	if err != nil {
		return errors.Wrapf(err, "deleting product: %v", err)
	}
	return nil
}
