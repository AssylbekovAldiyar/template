package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"
)

type Order struct {
	Id          int            `json:"id"`
	EmployeeID  int            `json:"employee_id"`
	TotalPrice  float64        `json:"total_price"`
	TotalPaid   float64        `json:"total_paid"`
	TotalReturn float64        `json:"total_return"`
	ReceiptID   string         `json:"receipt_id"`
	Products    []OrderProduct `json:"products"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type OrderModule struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (o OrderModule) Create(order *Order) error {
	query := `
			INSERT INTO orders (employee_id, total_price, total_paid, total_return, receipt_id, created_at, updated_at, products)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
				RETURNING id
			`
	// Serialize products slice to JSON
	productsJSON, err := json.Marshal(order.Products)
	if err != nil {
		return err
	}

	args := []interface{}{order.EmployeeID, order.TotalPrice, order.TotalPaid, order.TotalReturn, order.ReceiptID, order.CreatedAt, order.UpdatedAt,productsJSON}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return o.DB.QueryRowContext(ctx, query, args...).Scan(&order.Id)
}

func (o OrderModule) Get(id int) (*Order, error) {
    query := `
        SELECT * FROM orders 
        WHERE id = $1
    `
    var order Order
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    row := o.DB.QueryRowContext(ctx, query, id)
    var productsJSON []byte
    err := row.Scan(&order.Id, &order.EmployeeID, &order.TotalPrice, &order.TotalPaid,
        &order.TotalReturn, &order.ReceiptID, &order.CreatedAt, &order.UpdatedAt, &productsJSON)

    if err != nil {
        return nil, err
    }

    if err := json.Unmarshal(productsJSON, &order.Products); err != nil {
        return nil, err
    }

    return &order, nil
}


func (o OrderModule) GetAll() (*[]Order, error) {
    query := `SELECT * FROM orders`

    var orders []Order
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    rows, err := o.DB.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var ord Order
        var productsJSON []byte

        err := rows.Scan(&ord.Id, &ord.EmployeeID, &ord.TotalPrice, &ord.TotalPaid,
            &ord.TotalReturn, &ord.ReceiptID, &ord.CreatedAt, &ord.UpdatedAt, &productsJSON)
        if err != nil {
            return nil, err
        }

        if err := json.Unmarshal(productsJSON, &ord.Products); err != nil {
            return nil, err
        }

        orders = append(orders, ord)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return &orders, nil
}

func (o OrderModule) Update(id int, order *Order) error {
    query := `
        UPDATE orders
        SET employee_id = $1, total_price = $2, total_paid = $3, total_return = $4, receipt_id = $5, products = $6, updated_at = $7
        WHERE id = $8
        RETURNING updated_at
    `

    productsJSON, err := json.Marshal(order.Products)
    if err != nil {
        return err
    }

    args := []interface{}{order.EmployeeID, order.TotalPrice, order.TotalPaid, order.TotalReturn, order.ReceiptID, productsJSON, time.Now(), id}
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    return o.DB.QueryRowContext(ctx, query, args...).Scan(&order.UpdatedAt)
}

func (o OrderModule) Delete(id int) error {
	query := `
			DELETE FROM orders 
			WHERE id = $1
			`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := o.DB.ExecContext(ctx, query, id)
	return err
}
