package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBClient struct {
	Client *mongo.Client
	Ctx    context.Context
	Cancel context.CancelFunc
}

func InitDB() (*DBClient, error) {
	err := godotenv.Load(".env")

	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return nil, fmt.Errorf("MONGODB_URI env variable is not set")
	}

	// Create a context with a timeout for the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	// Configure MongoDB client options
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOpts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// connect to MongoDB
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}
	// ping the db
	if err := client.Ping(ctx, nil); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	bgCtx, bgCancel := context.WithCancel(context.Background())

	log.Println("Successfully connected to MongoDB")

	return &DBClient{
		Client: client,
		Ctx:    bgCtx,
		Cancel: bgCancel,
	}, nil

}

// GetDatabase returns a specific db from the client
func (db *DBClient) GetDatabase(name string) *mongo.Database {
	return db.Client.Database(name)
}

// GetCollection returns a specific collection from the db
func (db *DBClient) GetCollection(dbName string, collName string) *mongo.Collection {
	return db.Client.Database(dbName).Collection(collName)
}

// Close properly closes the MongoDB connection
func (db *DBClient) Close() error {

	// Cancel the bg context
	db.Cancel()

	// create a timeour context for disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// disconnect the client
	if err := db.Client.Disconnect(ctx); err != nil {
		return fmt.Errorf("error disconnecting from MongoDB: %v", err)
	}
	log.Println("Successfully disconnected from MongoDB")
	return nil

}
