package model

import (
	"database/sql"
	"errors"
	"log"
	"os"
)

var (
	// ErrRecordNotFound is returned when a record doesn't exist in database.
	ErrRecordNotFound = errors.New("record not found")

	// ErrEditConflict is returned when a there is a data race, and we have an edit conflict.
	ErrEditConflict = errors.New("edit conflict")
)

type Models struct {
	Employee    EmployeeModel
	Product     ProductModule
	Category    CategoryModule
	Order       OrderModule
	Tokens      TokenModel
	Permissions PermissionModel
}

func NewModels(db *sql.DB) Models {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return Models{
		Permissions: PermissionModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Employee: EmployeeModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Product: ProductModule{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Category: CategoryModule{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Order: OrderModule{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Tokens: TokenModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
	}
}
