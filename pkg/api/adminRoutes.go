package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/blackpanther26/mvc/pkg/controllers"
	"github.com/blackpanther26/mvc/pkg/middleware"
)

func AdminRoutes(r *mux.Router) {
	admin := r.PathPrefix("/admin").Subrouter()

	admin.Use(middleware.RequireAuth)
	admin.Use(middleware.NoCache)
    admin.Use(middleware.IsAdmin)

	admin.HandleFunc("/", controllers.AdminListBooks).Methods(http.MethodGet)
	admin.HandleFunc("/books/add", controllers.AdminAddBook).Methods(http.MethodGet, http.MethodPost)
	admin.HandleFunc("/books/{id}/edit", controllers.AdminEditBook).Methods(http.MethodGet, http.MethodPost)
	admin.HandleFunc("/books/{id}/delete", controllers.AdminDeleteBook).Methods(http.MethodPost)
	admin.HandleFunc("/transactions", controllers.AdminListTransactions).Methods(http.MethodGet)
    admin.HandleFunc("/transactions/{id}/approve", controllers.AdminApproveTransaction).Methods(http.MethodPost)
    admin.HandleFunc("/transactions/{id}/deny", controllers.AdminDenyTransaction).Methods(http.MethodPost)	
	admin.HandleFunc("/requests", controllers.AdminListAdminRequests).Methods(http.MethodGet)
	admin.HandleFunc("/requests/{id}/approve", controllers.AdminApproveAdminRequest).Methods(http.MethodPost)
	admin.HandleFunc("/requests/{id}/deny", controllers.AdminDenyAdminRequest).Methods(http.MethodPost)
}