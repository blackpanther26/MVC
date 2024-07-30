package api

import (
	"log"
	"net/http"
	"github.com/blackpanther26/mvc/pkg/config"
	"github.com/gorilla/mux"
)
func init() {
	config.LoadEnvVariables()
	config.ConnectToDb()
	config.SyncDatabase()
}

func Start() {
	r := mux.NewRouter()

	AuthRoutes(r)
	AdminRoutes(r)
	ClientRoutes(r)

	port := config.GetPort()
    log.Printf("Listening on port %s\n", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}