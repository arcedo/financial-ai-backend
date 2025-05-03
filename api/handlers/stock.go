package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/arcedo/financial-ai-backend/api/helpers"
	db "github.com/arcedo/financial-ai-backend/database"
	"github.com/arcedo/financial-ai-backend/types"
	"go.mongodb.org/mongo-driver/bson"
)

func GetAllStocks(w http.ResponseWriter, r *http.Request, store db.MongoStorage) error {
	stockCollection := store.Collection("stocks")
	cursor, err := stockCollection.Find(context.Background(), bson.D{})
	if err != nil {
		return fmt.Errorf("error retrieving stocks: %v", err)
	}
	defer cursor.Close(context.Background())

	var stocks []types.Stock
	for cursor.Next(context.Background()) {
		var stock types.Stock
		if err := cursor.Decode(&stock); err != nil {
			return fmt.Errorf("error decoding stocks: %v", err)
		}
		stocks = append(stocks, stock)
	}

	if err := cursor.Err(); err != nil {
		return fmt.Errorf("cursor error: %v", err)
	}

	helpers.WriteJSON(w, http.StatusOK, stocks, nil, "")
	return nil
}

func GetProducts(w http.ResponseWriter, r *http.Request, store db.MongoStorage) error {
	productsCollection := store.Collection("products")
	cursor, err := productsCollection.Find(context.Background(), bson.D{})
	if err != nil {
		return fmt.Errorf("error retrieving products: %v", err)
	}
	defer cursor.Close(context.Background())

	var products []types.Product
	for cursor.Next(context.Background()) {
		var product types.Product
		if err := cursor.Decode(&product); err != nil {
			return fmt.Errorf("error decoding products: %v", err)
		}
		products = append(products, product)
	}

	if err := cursor.Err(); err != nil {
		return fmt.Errorf("cursor error: %v", err)
	}

	helpers.WriteJSON(w, http.StatusOK, products, nil, "")
	return nil
}
