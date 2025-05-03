package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/arcedo/financial-ai-backend/api"
	"github.com/arcedo/financial-ai-backend/data"
	db "github.com/arcedo/financial-ai-backend/database"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoUri := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)
	mongoStorage, err := db.NewMongoStorage(mongoUri, "financial-ai")
	if err != nil {
		log.Fatalf("Mongo connection failed: %v", err)
	}
	defer mongoStorage.Close(context.Background())

	if err := mongoStorage.InitProducts(context.Background(), "products", data.Products); err != nil {
		log.Fatalf("Error initializing products: %v", err)
	}
	/*if err := mongoStorage.RemoveCollection(context.Background(), "products"); err != nil {
		log.Fatalf("Error removing products: %v", err)
	}*/

	/* We can't execute more queries in max 25 per day ;(

	ctx := context.Background()
	go func() {
		for {
			log.Println("Running stock sync...")
			if err := requests.SyncDailyStockData(ctx, *mongoStorage); err != nil {
				log.Printf("Sync error: %v", err)
			} else {
				log.Println("Stock sync completed successfully.")
			}

			time.Sleep(24 * time.Hour)
		}
	}()*/

	server := api.NewServer(os.Getenv("LISTEN_ADDRESS"), *mongoStorage)

	if err := server.Start(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
