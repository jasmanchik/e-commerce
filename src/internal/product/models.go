package product

import (
	"time"
)

type Product struct {
	ID          string    `json:"id" db:"product_id"`
	Name        string    `json:"name"`
	Cost        int       `json:"cost"`
	Quantity    int       `json:"quantity"`
	Sold        int       `json:"sold"`
	Revenue     int       `json:"revenue"`
	DateCreated time.Time `json:"date_created" db:"date_created"`
	DateUpdated time.Time `json:"date_updated" db:"date_updated"`
}

type NewProduct struct {
	Name     string `json:"name" validate:"required"`
	Cost     int    `json:"cost" validate:"required,gt=0"`
	Quantity int    `json:"quantity" validate:"required,gt=0"`
}

type Sale struct {
	ID          string    `db:"sale_id" json:"id"`
	ProductID   string    `db:"product_id" json:"product_id"`
	Quantity    int       `db:"quantity" json:"quantity"`
	Paid        int       `db:"paid" json:"paid"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
}

type NewSale struct {
	Quantity int `json:"quantity"`
	Paid     int `json:"paid"`
}
