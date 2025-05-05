package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pos-rs/pkg/pos/model"
	"strconv"

	"github.com/gorilla/mux"
)

func (app *Application) createOrder(w http.ResponseWriter, r *http.Request) {
	var newOrder model.Order

	err := json.NewDecoder((r.Body)).Decode(&newOrder)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalies Request Payload")
		return
	}

	err = app.Models.Order.Create(&newOrder)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	app.respondWithJSON(w, http.StatusCreated, newOrder)
}

func (app *Application) getOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	orderId, err := strconv.Atoi(param)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Order ID")
		return
	}

	Order, err := app.Models.Order.Get(orderId)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, err.Error())
		return 
	}

	app.respondWithJSON(w, http.StatusFound, Order)
}

func (app *Application) getAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := app.Models.Order.GetAll()
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	app.respondWithJSON(w, http.StatusFound, orders)
}

func (app *Application) addProductToOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	orderId, err := strconv.Atoi(param)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Order ID")
		return
	}

	var product model.OrderProduct
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	existingOrder, err := app.Models.Order.Get(orderId)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "Order Not Found")
		return
	}

	existingOrder.Products = append(existingOrder.Products, product)
	existingOrder.TotalPrice += float64(product.Price) * float64(product.Qty)

	err = app.Models.Order.Update(orderId, existingOrder)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	app.respondWithJSON(w, http.StatusOK, existingOrder)
}

func removeProduct(products []model.OrderProduct, productId int) []model.OrderProduct {
	var updatedProducts []model.OrderProduct
	for _, p := range products {
		if p.Id != productId {
			updatedProducts = append(updatedProducts, p)
		}
	}
	if updatedProducts != nil{
		return updatedProducts
	}
	return []model.OrderProduct{}
}

func (app *Application) removeProductFromOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderIDStr := vars["id"]
	productIDStr := vars["productId"]

	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Order ID")
		return
	}

	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	existingOrder, err := app.Models.Order.Get(orderID)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "Order Not Found")
		return
	}

	updatedProducts := removeProduct(existingOrder.Products, productID)
	updatedTotalPrice := calculateTotalPrice(updatedProducts)

	existingOrder.Products = updatedProducts
	existingOrder.TotalPrice = updatedTotalPrice

	err = app.Models.Order.Update(orderID, existingOrder)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	app.respondWithJSON(w, http.StatusOK, existingOrder)
}

func calculateTotalPrice(products []model.OrderProduct) float64 {
	totalPrice := 0.0
	for _, p := range products {
		totalPrice += float64(p.Price) * float64(p.Qty)
	}
	return totalPrice
}

func (app *Application) deleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	orderId, err := strconv.Atoi(param)
	fmt.Println(orderId)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid Order ID")
		return
	}

	err = app.Models.Order.Delete(orderId)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
