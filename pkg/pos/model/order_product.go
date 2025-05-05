package model

import (
	"time"
)

type OrderProduct struct {
	Id                int    `json:"id"`
	OrderId           string    `json:"order_id"`
	ProductId         string    `json:"product_id"`
	Qty               int       `json:"qty"`
	Price             int       `json:"price"`
	TotalNormalPrice  int       `json:"total_normal_price"`
	Product           Product 	`json:"product"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
