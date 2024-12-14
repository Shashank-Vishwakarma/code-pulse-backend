package main

import (
	"log"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/database"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/routes"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/config"
	"github.com/gin-gonic/gin"
)

func init() {
	// Read the config file
	err := config.NewEnv()
	if err != nil {
		log.Fatal(err)
	}

	// Connect to the database
	err = database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to the database")
}

func main() {
	r := gin.Default()

	// register middlewares

	// register routes
	routes.AuthRoutes(r)

	err := r.Run(":" + config.Config.PORT)
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
