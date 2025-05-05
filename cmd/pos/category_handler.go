package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pos-rs/pkg/pos/model"
	"strconv"

	"github.com/gorilla/mux"
)

func (app *Application) createCategory(w http.ResponseWriter, r *http.Request) {
	var newCategory model.Category

	err := json.NewDecoder(r.Body).Decode(&newCategory)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	err = app.Models.Category.Create(&newCategory)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	app.respondWithJSON(w, http.StatusCreated, newCategory)
}

func (app *Application) getCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["categoryId"]

	categoryId, err := strconv.Atoi(param)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Category ID")
		return
	}

	Category, err := app.Models.Category.Get(categoryId)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusFound, Category)
}

func (app *Application) getAllCategory(w http.ResponseWriter, r *http.Request) {
	Categorys, err := app.Models.Category.GetAll()
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println("Hello from getAllCategory")
	app.respondWithJSON(w, http.StatusFound, Categorys)
}

func (app *Application) updateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["categoryId"]

	categoryId, err := strconv.Atoi(param)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Category ID")
		return
	}

	var updatedCategory model.Category
	err = json.NewDecoder(r.Body).Decode(&updatedCategory)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	err = app.Models.Category.Update(categoryId, &updatedCategory)
	updatedCategory.Id = categoryId
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	app.respondWithJSON(w, http.StatusOK, updatedCategory)
}

func (app *Application) deleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["categoryId"]

	categoryId, err := strconv.Atoi(param)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Category ID")
		return
	}

	err = app.Models.Category.Delete(categoryId)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
