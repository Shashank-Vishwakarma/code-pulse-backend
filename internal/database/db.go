package database

import (
	"context"
	"fmt"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DBClient *mongo.Client

func Connect() error {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.Config.DATABASE_URL))
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	DBClient = client
	return nil
}
