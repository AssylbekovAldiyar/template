package model

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"log"
	"time"
	//"google.golang.org/protobuf/types/known/emptypb"
)

var AnonymousEmployee = &Employee{}

type Employee struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Surname     string    `json:"surname"`
	Password    string    `json:"password"`
	IsAdmin     bool      `json:"isAdmin"`
	Activated   bool      `json:"activated"`
	PhoneNumber string    `json:"phoneNumber"`
	Enrolled    time.Time `json:"enrolled"`
}

type EmployeeModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (e *Employee) IsAnonymous() bool {
	return e == AnonymousEmployee
	}

func (e EmployeeModel) Register(emp *Employee) error {
	query := `
			INSERT INTO employee (name, surname, password, is_admin, activated, phone_number, enrolled) 
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id, password
			`
	args := []interface{}{emp.Name, emp.Surname, emp.Password, emp.IsAdmin, emp.Activated, emp.PhoneNumber, emp.Enrolled}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return e.DB.QueryRowContext(ctx, query, args...).Scan(&emp.Id, &emp.Password)
}

func (m EmployeeModel) GetForToken(tokenScope, tokenPlaintext string) (*Employee, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	query := `
	SELECT employee.id, employee.name, employee.surname, employee.password, employee.activated, employee.is_admin
	FROM employee
	INNER JOIN tokens
	ON employee.id = tokens.user_id
	WHERE tokens.hash = $1
	AND tokens.scope = $2
	AND tokens.expiry > $3`

	args := []interface{}{tokenHash[:], tokenScope, time.Now()}
	var emp Employee
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&emp.Id,
		&emp.Name,
		&emp.Surname,
		&emp.Password,
		&emp.Activated,
		&emp.IsAdmin,
	)
	if err != nil {
		return nil, err
	}

	return &emp, nil
}

func (e EmployeeModel) Get(id int) (*Employee, error) {
	query := `SELECT * FROM employee where id = $1`

	var emp Employee
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := e.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&emp.Id, &emp.Name, &emp.Surname, &emp.Password, &emp.IsAdmin, &emp.Activated,  &emp.PhoneNumber, &emp.Enrolled)

	if err != nil {
		return nil, err
	}

	return &emp, nil
}

func (e EmployeeModel) GetAll() (*[]Employee, error) {
	query := `SELECT * from employee`

	var emp []Employee
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := e.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var employee Employee
		err := rows.Scan(&employee.Id, &employee.Name, &employee.Surname, &employee.Password, &employee.IsAdmin, &employee.PhoneNumber, &employee.Enrolled)
		if err != nil {
			return nil, err
		}
		emp = append(emp, employee)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &emp, nil
}

func (e EmployeeModel) Update(id int, emp *Employee) error {
	query := `
			UPDATE employee 
			SET name = $1, surname = $2, password = $3, is_admin = $4, activated = $5, phone_number = $6
			WHERE id = $7
			RETURNING id, password
	`
	args := []interface{}{emp.Name, emp.Surname, emp.Password, emp.IsAdmin, emp.Activated,emp.PhoneNumber,id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return e.DB.QueryRowContext(ctx, query, args...).Scan(&emp.Id, &emp.Password)
}

func (e EmployeeModel) Delete(id int) error {
	query := `DELETE FROM employee WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := e.DB.ExecContext(ctx, query, id)
	return err
}
