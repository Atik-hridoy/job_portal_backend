package configs

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB configuration
const (
	MongoDBURI     = "mongodb://localhost:27017"
	MongoDBName    = "jobportal"
	MongoDBTimeout = 10 * time.Second
)

// DB holds the MongoDB client and database instance
type DB struct {
	Client *mongo.Client
	DB     *mongo.Database
}

// ConnectDB creates a new MongoDB connection
func ConnectDB() (*DB, error) {
	// Set client options
	clientOptions := options.Client().ApplyURI(MongoDBURI)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), MongoDBTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB!")

	return &DB{
		Client: client,
		DB:     client.Database(MongoDBName),
	}, nil
}

// Close disconnects from MongoDB
func (db *DB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), MongoDBTimeout)
	defer cancel()
	return db.Client.Disconnect(ctx)
}
