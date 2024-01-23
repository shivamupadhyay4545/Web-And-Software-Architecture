package main

import (
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"

	"github.com/shivamupadhyay4545/Web-And-Software-Architecture/service/api/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := httprouter.New()

	// Apply custom CORS middleware
	handler := cors.AllowAll().Handler(router)

	// Register custom middleware for logging
	handler = logRequests(handler)

	// Define routes
	routes.PhotoRoutes(router)
	routes.UserRoutes(router)

	// Start the server
	log.Printf("Server is listening on port %s...\n", port)
	err := http.ListenAndServe(":"+port, handler)
	if err != nil {
		log.Fatal(err)
	}
}

// logRequests is a middleware function that logs information about incoming requests
func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request details
		log.Printf("Request: %s %s %s\n", r.Method, r.URL.Path, r.RemoteAddr)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
