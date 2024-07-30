package main

import (
	"fmt"

	"github.com/blackpanther26/mvc/pkg/api"
)


func main() {
	fmt.Println("Started the API server")
	api.Start()
}

// go build -o mvc ./cmd/main.go