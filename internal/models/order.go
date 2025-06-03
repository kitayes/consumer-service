package models

import "time"

type Order struct {
	ID          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id"`
	ProductName string    `json:"product_name" db:"product_name"`
	Quantity    int       `json:"quantity" db:"quantity"`
	Price       float64   `json:"price" db:"price"`
	TotalPrice  float64   `json:"total_price" db:"total_price"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
