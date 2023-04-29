package schema

import (
	"github.com/jmoiron/sqlx"
)

const seed = `
INSERT INTO products(product_id, name, cost, quantity, date_created, date_updated)
VALUES ('4eabde0b-3331-4927-9091-701f829b0262', 'Comic Books', 50, 42, NOW(), NOW()),
       ('b0de7d30-42e4-4ee2-8f1a-a382be080c32', 'McDonalds Toys', 75, 120, NOW(), NOW()),
       ('2378aa21-db61-4d71-b7c8-3ee573df000a', 'Big Wheels', 500, 2, NOW(), NOW())
ON CONFLICT DO NOTHING;`

func Seed(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec(seed); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
	}

	return tx.Commit()
}
