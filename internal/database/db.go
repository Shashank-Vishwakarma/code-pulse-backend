package database

import (
	"context"
	"fmt"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/config"
	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/constants"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DBClient *mongo.Client
var UserCollection *mongo.Collection

func Connect() error {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.Config.DATABASE_URL))
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	DBClient = client

	// Initialize collections
	UserCollection = getOrCreateCollection(constants.USER_COLLECTION)

	return nil
}

func getOrCreateCollection(collectionName string) *mongo.Collection {
	return DBClient.Database(config.Config.DATABASE_NAME).Collection(collectionName)
}
