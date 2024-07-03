package api

import (
	"log"
	"net/http"
	"os"
	"github.com/blackpanther26/mvc/pkg/config"
	"github.com/gorilla/mux"
	// "github.com/blackpanther26/mvc/pkg/routes"
)
func init() {
	config.LoadEnvVariables()
	config.ConnectToDb()
	config.SyncDatabase()
}

func Start() {
	r := mux.NewRouter()

	AuthRoutes(r)
	// routes.AdminRoutes(r)
	ClientRoutes(r)
	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set")
	}

	log.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}