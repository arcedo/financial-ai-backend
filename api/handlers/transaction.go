package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/arcedo/financial-ai-backend/api/helpers"
	db "github.com/arcedo/financial-ai-backend/database"
	"github.com/arcedo/financial-ai-backend/types"
	"github.com/arcedo/financial-ai-backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateTransaction(w http.ResponseWriter, r *http.Request, store db.MongoStorage) error {
	// Parse the request body
	var transaction types.TransactionPublic
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	transaction.Type = utils.SanitizeString(transaction.Type)

	// Retrieve user ID from context
	userID, ok := r.Context().Value("userID").(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("unable to retrieve user ID from context")
	}

	// Validate transaction data
	if err := types.ValidateTransaction(transaction); err != nil {
		return err
	}

	// Check if product exists if type is buy/sell
	if transaction.Type == "buy" || transaction.Type == "sell" {
		productsColl := store.Collection("products")
		var product types.Product
		if err := productsColl.FindOne(r.Context(), bson.M{"symbol": transaction.Symbol}).Decode(&product); err != nil {
			if err == mongo.ErrNoDocuments {
				return fmt.Errorf("no product found with that symbol")
			}
			return err
		}
	} else {
		transaction.Symbol = ""
	}

	// Parse date
	dateFormated, err := time.Parse("2006-01-02", transaction.Date)
	if err != nil {
		return fmt.Errorf("invalid date format, expected YYYY-MM-DD")
	}

	// Prepare the transaction to insert
	newTransaction := types.NewTransaction{
		UserID: userID,
		Symbol: transaction.Symbol,
		Type:   transaction.Type,
		Amount: transaction.Amount,
		Date:   dateFormated.Format("2006-01-02"),
	}

	// Insert into database
	transactionsColl := store.Collection("transactions")
	_, err = transactionsColl.InsertOne(r.Context(), newTransaction)
	if err != nil {
		return fmt.Errorf("failed to insert transaction: %v", err)
	}

	// Return success
	helpers.WriteJSON(w, http.StatusCreated, transaction, nil, "transaction created successfully")
	return nil
}

func GetTransactions(w http.ResponseWriter, r *http.Request, store db.MongoStorage) error {
	// Retrieve user ID from context
	userID, ok := r.Context().Value("userID").(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("unable to retrieve user ID from context")
	}

	// Fetch transactions for the user (fix: don't use .Hex())
	transactionsColl := store.Collection("transactions")
	cursor, err := transactionsColl.Find(r.Context(), bson.M{"user_id": userID})
	if err != nil {
		return fmt.Errorf("failed to fetch transactions: %v", err)
	}
	defer cursor.Close(r.Context())

	var transactions []types.TransactionPublic
	if err := cursor.All(r.Context(), &transactions); err != nil {
		return fmt.Errorf("failed to decode transactions: %v", err)
	}

	helpers.WriteJSON(w, http.StatusOK, transactions, nil, "")
	return nil
}

func GetAllTransactions(w http.ResponseWriter, r *http.Request, store db.MongoStorage) error {
	// Fetch all transactions
	transactionsColl := store.Collection("transactions")
	cursor, err := transactionsColl.Find(r.Context(), bson.D{})
	if err != nil {
		return fmt.Errorf("failed to fetch transactions: %v", err)
	}
	defer cursor.Close(r.Context())

	var transactions []types.Transaction
	if err := cursor.All(r.Context(), &transactions); err != nil {
		return fmt.Errorf("failed to decode transactions: %v", err)
	}

	if len(transactions) == 0 {
		helpers.WriteJSON(w, http.StatusOK, []types.TransactionPublic{}, nil, "No transactions found")
		return nil
	}

	helpers.WriteJSON(w, http.StatusOK, transactions, nil, "")
	return nil
}
