package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
	routes "github.com/shivamupadhyay4545/Web-And-Software-Architecture/service/api/routes"
)

func main() {

	if 1 == 1 {
		fmt.Println("done")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	// Apply CORS middleware with default settings
	handler := cors.Default().Handler(mux)

	// Define routes
	routes.PhotoRoutes(mux)
	routes.UserRoutes(mux)

	// Start the server
	log.Printf("Server is listening on port %s...\n", port)
	err := http.ListenAndServe(":"+port, handler)
	if err != nil {
		log.Fatal(err)
	}
}
