package main

import (
	"log"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/database"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/queue"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/routes"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/internal/services"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/config"
	"github.com/gin-contrib/cors"
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

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	corsConfig.AllowCredentials = true

	r.Use(cors.New(corsConfig))

	// register routes
	routes.AuthRoutes(r)
	routes.QuestionRoutes(r)

	// rabbitmq setup
	err := queue.InitializeRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}

	// start the consumer
	err = queue.StartConsumer()
	if err != nil {
		log.Fatalf("Failed to start the consumer: %v", err)
	}

	// Connect to redis
	services.InitializeRedis()

	// start the server
	err = r.Run(":" + config.Config.PORT)
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
