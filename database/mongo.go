package db

import (
	"context"
	"fmt"
	"time"

	"github.com/arcedo/financial-ai-backend/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoStorage wraps a MongoDB client and database
type MongoStorage struct {
	client   *mongo.Client
	database *mongo.Database
}

// NewMongoStorage creates a new MongoDB connection and returns a storage instance
func NewMongoStorage(uri, dbName string) (*MongoStorage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping to confirm the connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := client.Database(dbName)
	return &MongoStorage{client: client, database: db}, nil
}

// Collection returns a reference to a MongoDB collection
func (m *MongoStorage) Collection(name string) *mongo.Collection {
	return m.database.Collection(name)
}

// Close gracefully disconnects the MongoDB client
func (m *MongoStorage) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

func (m *MongoStorage) InitProducts(ctx context.Context, collection string, products []types.Product) error {
	col := m.database.Collection(collection)

	count, err := col.CountDocuments(ctx, bson.D{})
	if err != nil {
		return fmt.Errorf("failed to count documents: %w", err)
	}
	if count > 0 {
		return nil
	}

	docs := make([]interface{}, len(products))
	for i, product := range products {
		docs[i] = product
	}

	_, err = col.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("failed to insert products: %w", err)
	}

	return nil
}

func (m *MongoStorage) RemoveCollection(ctx context.Context, collection string) error {
	col := m.database.Collection(collection)

	_, err := col.DeleteMany(ctx, bson.D{})
	if err != nil {
		return fmt.Errorf("failed to remove all products: %w", err)
	}

	return nil
}
