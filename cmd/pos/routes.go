package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (app *Application) routes() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/healthcheck", app.healthcheckHandler).Methods("GET")

	v1 := r.PathPrefix("/api/v1").Subrouter()
	fmt.Println("Running")

	v1.HandleFunc("/tokens/authentication", app.createAuthenticationTokenHandler).Methods("POST")
	v1.HandleFunc("/employees", app.getAllEmployee).Methods("GET")
	v1.HandleFunc("/employees/activated", app.activateUserHandler).Methods("PUT")
	v1.HandleFunc("/employees/{id}", app.getEmployee).Methods("GET")
	v1.HandleFunc("/employees", app.registerEmployee).Methods("POST")
	v1.HandleFunc("/employees/{id}", app.updateEmployee).Methods("PUT")
	v1.HandleFunc("/employees/{id}", app.requireActivatedUser(app.deleteEmployee)).Methods("DELETE")

	v1.HandleFunc("/categories", app.getAllCategory).Methods("GET")
	v1.HandleFunc("/categories/{categoryId}", app.getCategory).Methods("GET")
	v1.HandleFunc("/categories", app.createCategory).Methods("POST")
	v1.HandleFunc("/categories/{categoryId}", app.updateCategory).Methods("PUT")
	v1.HandleFunc("/categories/{categoryId}", app.deleteCategory).Methods("DELETE")

	v1.HandleFunc("/products", app.getAllProduct).Methods("GET")
	v1.HandleFunc("/products/{productId}", app.getProduct).Methods("GET")
	v1.HandleFunc("/products", app.createProduct).Methods("POST")
	v1.HandleFunc("/products/{productId}", app.updateProduct).Methods("PUT")
	v1.HandleFunc("/products/{productId}", app.requirePermission("products:write",app.deleteProduct)).Methods("DELETE")

	v1.HandleFunc("/orders", app.getAllOrders).Methods("GET")
	v1.HandleFunc("/orders/{id}", app.getOrder).Methods("GET")
	v1.HandleFunc("/orders", app.createOrder).Methods("POST")
	v1.HandleFunc("/orders/{id}/products", app.addProductToOrder).Methods("PUT")
	v1.HandleFunc("/orders/{id}/products/{productId}", app.removeProductFromOrder).Methods("PUT")
	v1.HandleFunc("/orders/{id}", app.deleteOrder).Methods("DELETE")

	return app.recoverPanic(app.rateLimit(app.authenticate(r)))
}
