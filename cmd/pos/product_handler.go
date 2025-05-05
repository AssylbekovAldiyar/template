package main

import (
	"encoding/json"

	"net/http"
	"pos-rs/pkg/pos/model"
	"pos-rs/pkg/pos/validator"
	"strconv"

	"github.com/gorilla/mux"
)

func (app *Application) createProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct model.Product

	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	err = app.Models.Product.Create(&newProduct)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	app.respondWithJSON(w, http.StatusCreated, newProduct)
}

func (app *Application) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["productId"]

	productId, err := strconv.Atoi(param)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	Product, err := app.Models.Product.Get(productId)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusFound, Product)
}

func (app *Application) getAllProduct(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string
		Cateogry int
		model.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Name = app.readString(qs, "name", "")
	input.Cateogry = app.readInt(qs, "category", 1, v)

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "price", "-id", "-name", "-price"}

	if !v.Valid() {
		app.respondWithError(w, http.StatusForbidden, "Failed Validation")
	}

	products, metadata, err := app.Models.Product.GetAll(input.Name, input.Cateogry, input.Filters)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	app.respondWithJSON(w, http.StatusFound, envelope{"products": products, "metadata": metadata})
}

func (app *Application) updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["productId"]

	productId, err := strconv.Atoi(param)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	var updatedProduct model.Product
	err = json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	err = app.Models.Product.Update(productId, &updatedProduct)
	updatedProduct.Id = productId
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	app.respondWithJSON(w, http.StatusOK, updatedProduct)
}

func (app *Application) deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["productId"]

	productId, err := strconv.Atoi(param)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}
	err = app.Models.Product.Delete(productId)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
