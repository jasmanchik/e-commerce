package product

import (
	"time"
)

type Product struct {
	ID          string    `json:"id" db:"product_id"`
	Name        string    `json:"name"`
	Cost        int       `json:"cost"`
	Quantity    int       `json:"quantity"`
	DateCreated time.Time `json:"date_created" db:"date_created"`
	DateUpdated time.Time `json:"date_updated" db:"date_updated"`
}

type NewProduct struct {
	Name     string `json:"name" validate:"required"`
	Cost     int    `json:"cost" validate:"required,gt=0"`
	Quantity int    `json:"quantity" validate:"required,gt=0"`
}
