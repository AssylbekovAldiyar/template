package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Category struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CategoryModule struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (c CategoryModule) Create(category *Category) error {
	query := `
			INSERT INTO categories (name)
			VALUES ($1)
			RETURNING id
			`
	args := []interface{}{category.Name}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.DB.QueryRowContext(ctx, query, args...).Scan(&category.Id)
}

func (c CategoryModule) GetAll() (*[]Category, error) {
	query := `SELECT * from categories`

	var categories []Category
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := c.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ctg Category
		err := rows.Scan(&ctg.Id, &ctg.Name, &ctg.CreatedAt, &ctg.UpdatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, ctg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &categories, nil
}

func (c CategoryModule) Get(id int) (*Category, error) {
	query := `
			SELECT * FROM categories
			WHERE id = $1
			`

	var category Category
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := c.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (c CategoryModule) Update(id int, category *Category) error {
	query := `
			UPDATE categories 
			SET name = $1, updated_at = CURRENT_TIMESTAMP 
			WHERE id = $2
			`
	args := []interface{}{category.Name, id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := c.DB.ExecContext(ctx, query, args...)
	return err
}

func (c CategoryModule) Delete(id int) error {
	query := `
			DELETE FROM categories
			WHERE id = $1
			`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := c.DB.ExecContext(ctx, query, id)
	return err
}
