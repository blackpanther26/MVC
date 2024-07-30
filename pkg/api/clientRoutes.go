package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/blackpanther26/mvc/pkg/controllers"
	"github.com/blackpanther26/mvc/pkg/middleware"
)

func ClientRoutes(r *mux.Router) {
	client := r.PathPrefix("/client").Subrouter()

	client.Use(middleware.RequireAuth)
	client.Use(middleware.NoCache)
    client.Use(middleware.IsNotAdmin)

	client.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        controllers.ListBooks(w, r, false) 
    }).Methods(http.MethodGet)
	client.HandleFunc("/books/{id}/checkout", controllers.CheckoutBook).Methods(http.MethodPost)
	client.HandleFunc("/books/{id}/checkin", controllers.CheckinBook).Methods(http.MethodPost)
	client.HandleFunc("/history", controllers.UserHistory).Methods(http.MethodGet)
	client.HandleFunc("/requestAdmin", controllers.RequestAdmin).Methods(http.MethodPost)
	client.HandleFunc("/search", controllers.SearchBooks).Methods(http.MethodGet)
}