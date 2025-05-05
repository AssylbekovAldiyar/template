package main

import (
	"encoding/json"
	"fmt"

	// "go/token"
	"net/http"
	"pos-rs/pkg/pos/model"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func (app *Application) registerEmployee(w http.ResponseWriter, r *http.Request) {
	var newEmployee model.Employee

	err := json.NewDecoder(r.Body).Decode(&newEmployee)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newEmployee.Password), bcrypt.DefaultCost)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	newEmployee.Password = string(hashedPassword)
	newEmployee.Activated = false
	newEmployee.Enrolled = time.Now()



	err = app.Models.Employee.Register(&newEmployee)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	fmt.Println(newEmployee.Id)	
	err = app.Models.Permissions.AddForUser(newEmployee.Id, "products:read")
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	token, err := app.Models.Tokens.New(newEmployee.Id, 3*24*time.Hour, model.ScopeActivision)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	app.respondWithJSON(w, http.StatusCreated, token)
}

func (app *Application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TokenPlaintext string `json:"token"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println(input.TokenPlaintext)

	employee, err := app.Models.Employee.GetForToken(model.ScopeActivision, input.TokenPlaintext)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	employee.Activated = true

	err = app.Models.Employee.Update(employee.Id, employee)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = app.Models.Tokens.DeleteAllForUser(model.ScopeActivision, employee.Id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	app.respondWithJSON(w, http.StatusOK, envelope{"employee": employee})
}

func (app *Application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Id       int    `json:"id"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	employee, err := app.Models.Employee.Get(input.Id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(employee.Password), []byte(input.Password))

	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := app.Models.Tokens.New(employee.Id, 24*time.Hour, model.ScopeAuthentication)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	app.respondWithJSON(w, http.StatusCreated, envelope{"authentication_token": token})

}

func (app *Application) logInEmployee(w http.ResponseWriter, r *http.Request) {
	var logInRequest struct {
		Id       int    `json:"id"`
		Password string `string:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&logInRequest)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	Employee, err := app.Models.Employee.Get(logInRequest.Id)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(logInRequest.Password), 14)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(Employee.Password), hashedPassword)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	app.respondWithJSON(w, http.StatusOK, logInRequest)
}

func (app *Application) updateEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	employeeId, err := strconv.Atoi(param)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Employee ID")
		return
	}

	var updatedEmployee model.Employee
	err = json.NewDecoder(r.Body).Decode(&updatedEmployee)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	err = app.Models.Employee.Update(employeeId, &updatedEmployee)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, updatedEmployee)
}

func (app *Application) getAllEmployee(w http.ResponseWriter, r *http.Request) {
	employees, err := app.Models.Employee.GetAll()
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusFound, employees)
}

func (app *Application) getEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	employeeId, err := strconv.Atoi(param)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Employee ID")
		return
	}

	Employee, err := app.Models.Employee.Get(employeeId)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Employee Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusFound, Employee)
}

func (app *Application) deleteEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	employeeId, err := strconv.Atoi(param)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Employee ID")
		return
	}

	err = app.Models.Employee.Delete(employeeId)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
