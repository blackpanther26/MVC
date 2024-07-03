package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/blackpanther26/mvc/pkg/controllers"
	"github.com/blackpanther26/mvc/pkg/middleware"
)

func ClientRoutes(r *mux.Router) {
	client := r.PathPrefix("/client").Subrouter()

	client.Use(middleware.RequireAuth)
    client.Use(middleware.IsNotAdmin)

	client.HandleFunc("/books", controllers.ListBooks).Methods(http.MethodGet)
	client.HandleFunc("/books/{id}/checkout", controllers.CheckoutBook).Methods(http.MethodPost)
	client.HandleFunc("/books/{id}/checkin", controllers.CheckinBook).Methods(http.MethodPost)
	client.HandleFunc("/history", controllers.UserHistory).Methods(http.MethodGet)
}