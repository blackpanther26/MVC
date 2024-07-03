package routes

import (
	"github.com/blackpanther26/mvc/pkg/controllers"
	"github.com/gorilla/mux"
)

func AuthRoutes(r *mux.Router) {
	r.HandleFunc("/signup", controllers.SignupPageHandler).Methods("GET")
	r.HandleFunc("/signup", controllers.SignupHandler).Methods("POST")
	r.HandleFunc("/login", controllers.LoginPageHandler).Methods("GET")
	r.HandleFunc("/login", controllers.LoginHandler).Methods("POST")
}