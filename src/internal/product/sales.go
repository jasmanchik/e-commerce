package product

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

func AddSale(ctx context.Context, db *sqlx.DB, ns NewSale, productID string, now time.Time) (*Sale, error) {
	s := Sale{
		ID:          uuid.New().String(),
		ProductID:   productID,
		Quantity:    ns.Quantity,
		Paid:        ns.Paid,
		DateCreated: now.UTC(),
	}

	const q = `INSERT INTO sales (sale_id, product_id, quantity, paid, date_created) VALUES ($1,$2,$3,$4,$5)`

	_, err := db.ExecContext(ctx, q, s.ID, s.ProductID, s.Quantity, s.Paid, s.DateCreated)
	if err != nil {
		return nil, errors.Wrap(err, "inserting sale")
	}

	return &s, nil
}

func ListSales(ctx context.Context, db *sqlx.DB, productID string) ([]Sale, error) {

	var sl []Sale
	const q = `SELECT * from sales where product_id = $1`
	err := db.SelectContext(ctx, &sl, q, productID)
	if err != nil {
		return nil, errors.Wrap(err, "select sale")
	}

	return sl, nil

}
