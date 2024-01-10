package main

import (
	"os"

	routes "github.com/shivamupadhyay4545/Web-And-Software-Architecture/service/api/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	router := gin.New()
	router.Use(gin.Logger())

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.MaxAge = 1
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"} // Add other methods as needed
	router.Use(cors.New(config))

	routes.PhotoRoutes(router)
	routes.UserRoutes(router)

	router.Run(":" + port)

}
